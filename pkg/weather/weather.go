// Package weather provides temperature lookup for locations and date ranges
// using the free Open-Meteo API (no API key required).
package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// httpClient is used for all outbound requests so we can enforce a timeout.
var httpClient = &http.Client{Timeout: 10 * time.Second}

// Base URLs for Open-Meteo APIs. Tests override these to point at a mock server.
var (
	geocodingBaseURL = "https://geocoding-api.open-meteo.com"
	forecastBaseURL  = "https://api.open-meteo.com"
	archiveBaseURL   = "https://archive-api.open-meteo.com"
)

// Result holds the weather lookup results for a location and date range.
type Result struct {
	Location       string  `json:"location"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Nights         int     `json:"nights"`
	TemperatureMin int     `json:"temperatureMin"`
	TemperatureMax int     `json:"temperatureMax"`
	Source         string  `json:"source"` // "forecast" or "typical"
}

func localDate(t time.Time) time.Time {
	y, m, d := t.In(time.Local).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// geocodingResponse is the Open-Meteo geocoding API response.
type geocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
		Admin1    string  `json:"admin1"`
	} `json:"results"`
}

// dailyWeatherResponse is the Open-Meteo daily weather API response.
type dailyWeatherResponse struct {
	Daily struct {
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

// Lookup geocodes the location and fetches weather data for the given date
// range. Dates must be today or in the future. If the dates are within the
// 16-day forecast window, it uses the forecast API. If the dates are further
// out, it looks up historical data from the same dates in the previous year
// to get "typical" temperatures.
func Lookup(location string, startDate, endDate time.Time) (*Result, error) {
	startDate = localDate(startDate)
	endDate = localDate(endDate)

	// Reject past dates using the server's local calendar date.
	today := localDate(time.Now())
	if startDate.Before(today) {
		return nil, fmt.Errorf("start date must be today or in the future")
	}
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	// Step 1: Geocode the location
	lat, lon, resolvedName, err := geocode(location)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Calculate nights (endDate - startDate)
	nights := int(endDate.Sub(startDate).Hours()/24 + 0.5)
	if nights < 1 {
		nights = 1
	}

	// Step 2: Determine which API to use based on dates.
	// The Open-Meteo forecast API covers up to 16 days ahead, inclusive.
	forecastLimit := today.AddDate(0, 0, 16)

	var tempMin, tempMax float64
	var source string

	if !startDate.After(forecastLimit) && !endDate.After(forecastLimit) {
		// All dates within forecast window → use forecast API
		tempMin, tempMax, err = fetchForecast(lat, lon, startDate, endDate)
		source = "forecast"
	} else {
		// Dates too far in future → use same dates from previous year
		typicalStart := startDate.AddDate(-1, 0, 0)
		typicalEnd := endDate.AddDate(-1, 0, 0)
		tempMin, tempMax, err = fetchHistorical(lat, lon, typicalStart, typicalEnd)
		source = "typical"
	}

	if err != nil {
		return nil, fmt.Errorf("weather fetch failed: %w", err)
	}

	return &Result{
		Location:       resolvedName,
		Latitude:       lat,
		Longitude:      lon,
		Nights:         nights,
		TemperatureMin: int(math.Round(tempMin)),
		TemperatureMax: int(math.Round(tempMax)),
		Source:         source,
	}, nil
}

// LookupByCoords fetches weather data for known coordinates, skipping
// geocoding entirely. Used when the frontend already resolved the location
// via autocomplete.
func LookupByCoords(lat, lon float64, displayName string, startDate, endDate time.Time) (*Result, error) {
	startDate = localDate(startDate)
	endDate = localDate(endDate)

	today := localDate(time.Now())
	if startDate.Before(today) {
		return nil, fmt.Errorf("start date must be today or in the future")
	}
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	nights := int(endDate.Sub(startDate).Hours()/24 + 0.5)
	if nights < 1 {
		nights = 1
	}

	forecastLimit := today.AddDate(0, 0, 16)

	var tempMin, tempMax float64
	var source string
	var err error

	if !startDate.After(forecastLimit) && !endDate.After(forecastLimit) {
		tempMin, tempMax, err = fetchForecast(lat, lon, startDate, endDate)
		source = "forecast"
	} else {
		typicalStart := startDate.AddDate(-1, 0, 0)
		typicalEnd := endDate.AddDate(-1, 0, 0)
		tempMin, tempMax, err = fetchHistorical(lat, lon, typicalStart, typicalEnd)
		source = "typical"
	}

	if err != nil {
		return nil, fmt.Errorf("weather fetch failed: %w", err)
	}

	return &Result{
		Location:       displayName,
		Latitude:       lat,
		Longitude:      lon,
		Nights:         nights,
		TemperatureMin: int(math.Round(tempMin)),
		TemperatureMax: int(math.Round(tempMax)),
		Source:         source,
	}, nil
}

// (state/province) names for the US and Canada.
var admin1Abbreviations = map[string]string{
	// US States
	"al": "alabama", "ak": "alaska", "az": "arizona", "ar": "arkansas",
	"ca": "california", "co": "colorado", "ct": "connecticut", "de": "delaware",
	"fl": "florida", "ga": "georgia", "hi": "hawaii", "id": "idaho",
	"il": "illinois", "in": "indiana", "ia": "iowa", "ks": "kansas",
	"ky": "kentucky", "la": "louisiana", "me": "maine", "md": "maryland",
	"ma": "massachusetts", "mi": "michigan", "mn": "minnesota", "ms": "mississippi",
	"mo": "missouri", "mt": "montana", "ne": "nebraska", "nv": "nevada",
	"nh": "new hampshire", "nj": "new jersey", "nm": "new mexico", "ny": "new york",
	"nc": "north carolina", "nd": "north dakota", "oh": "ohio", "ok": "oklahoma",
	"or": "oregon", "pa": "pennsylvania", "ri": "rhode island", "sc": "south carolina",
	"sd": "south dakota", "tn": "tennessee", "tx": "texas", "ut": "utah",
	"vt": "vermont", "va": "virginia", "wa": "washington", "wv": "west virginia",
	"wi": "wisconsin", "wy": "wyoming", "dc": "district of columbia",
	// Canadian Provinces & Territories
	"ab": "alberta", "bc": "british columbia", "mb": "manitoba",
	"nb": "new brunswick", "nl": "newfoundland and labrador",
	"ns": "nova scotia", "nt": "northwest territories",
	"nu": "nunavut", "on": "ontario", "pe": "prince edward island",
	"qc": "quebec", "sk": "saskatchewan", "yt": "yukon",
}

// matchesAdmin1Filter checks whether an admin1 or country value matches the
// user's filter string. It handles both prefix matching on the full name and
// exact abbreviation matching (e.g. "ny" → "new york").
func matchesAdmin1Filter(admin1, country, filter string) bool {
	lowerAdmin1 := strings.ToLower(admin1)
	lowerCountry := strings.ToLower(country)

	// Direct prefix match on full name
	if strings.HasPrefix(lowerAdmin1, filter) || strings.HasPrefix(lowerCountry, filter) {
		return true
	}

	// Abbreviation match: expand the filter if it's a known abbreviation
	if expanded, ok := admin1Abbreviations[filter]; ok {
		if strings.EqualFold(lowerAdmin1, expanded) {
			return true
		}
	}

	return false
}

// GeocodeSuggestion represents a single geocoding result for autocomplete.
type GeocodeSuggestion struct {
	Name      string  `json:"name"`
	Admin1    string  `json:"admin1,omitempty"`
	Country   string  `json:"country,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Display   string  `json:"display"`
}

// GeocodeSuggestions returns up to 5 location suggestions for the given query
// using the Open-Meteo Geocoding API. If the query contains a comma, only the
// part before the comma is sent to the API as the place name, and the remainder
// is used to filter results by admin1 (state/province) or country.
func GeocodeSuggestions(query string) ([]GeocodeSuggestion, error) {
	searchName := query
	var filter string
	if idx := strings.Index(query, ","); idx >= 0 {
		searchName = strings.TrimSpace(query[:idx])
		filter = strings.ToLower(strings.TrimSpace(query[idx+1:]))
	}

	if searchName == "" {
		return nil, nil
	}

	// Request more results when filtering so we have enough after narrowing.
	count := 5
	if filter != "" {
		count = 20
	}

	u := fmt.Sprintf("%s/v1/search?name=%s&count=%d&language=en&format=json",
		geocodingBaseURL, url.QueryEscape(searchName), count)

	resp, err := httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding returned status %d", resp.StatusCode)
	}

	var geo geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return nil, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	suggestions := make([]GeocodeSuggestion, 0, 5)
	for _, r := range geo.Results {
		display := r.Name
		if r.Admin1 != "" {
			display += ", " + r.Admin1
		}
		if r.Country != "" {
			display += ", " + r.Country
		}

		// If the user typed something after a comma, filter results
		// by matching against admin1 (state/province) or country.
		if filter != "" {
			if !matchesAdmin1Filter(r.Admin1, r.Country, filter) {
				continue
			}
		}

		suggestions = append(suggestions, GeocodeSuggestion{
			Name:      r.Name,
			Admin1:    r.Admin1,
			Country:   r.Country,
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
			Display:   display,
		})

		if len(suggestions) >= 5 {
			break
		}
	}

	return suggestions, nil
}

// geocode uses the Open-Meteo Geocoding API to resolve a location name to
// coordinates. If location contains commas (e.g. "Chatham, New York, United
// States"), only the part before the first comma is sent as the search name
// and the rest is used to pick the best match.
func geocode(location string) (lat, lon float64, name string, err error) {
	searchName := location
	var filterParts []string
	if idx := strings.Index(location, ","); idx >= 0 {
		searchName = strings.TrimSpace(location[:idx])
		for _, part := range strings.Split(location[idx+1:], ",") {
			if p := strings.TrimSpace(part); p != "" {
				filterParts = append(filterParts, strings.ToLower(p))
			}
		}
	}

	count := 1
	if len(filterParts) > 0 {
		count = 20
	}

	u := fmt.Sprintf("%s/v1/search?name=%s&count=%d&language=en&format=json",
		geocodingBaseURL, url.QueryEscape(searchName), count)

	resp, err := httpClient.Get(u)
	if err != nil {
		return 0, 0, "", fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, "", fmt.Errorf("geocoding returned status %d", resp.StatusCode)
	}

	var geo geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return 0, 0, "", fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	if len(geo.Results) == 0 {
		return 0, 0, "", fmt.Errorf("location %q not found", location)
	}

	// If the caller passed extra qualifiers (state, country), find the
	// first result that matches all of them.
	r := geo.Results[0]
	if len(filterParts) > 0 {
		found := false
		for _, candidate := range geo.Results {
			fields := strings.ToLower(candidate.Admin1 + " " + candidate.Country)
			allMatch := true
			for _, fp := range filterParts {
				if !matchesAdmin1Filter(candidate.Admin1, candidate.Country, fp) &&
					!strings.Contains(fields, fp) {
					allMatch = false
					break
				}
			}
			if allMatch {
				r = candidate
				found = true
				break
			}
		}
		if !found {
			return 0, 0, "", fmt.Errorf("location %q not found", location)
		}
	}

	resolvedName := r.Name
	if r.Admin1 != "" {
		resolvedName += ", " + r.Admin1
	}
	if r.Country != "" {
		resolvedName += ", " + r.Country
	}

	return r.Latitude, r.Longitude, resolvedName, nil
}

// fetchForecast uses the Open-Meteo Forecast API (up to 16 days ahead).
func fetchForecast(lat, lon float64, start, end time.Time) (minTemp, maxTemp float64, err error) {
	u := fmt.Sprintf(
		"%s/v1/forecast?latitude=%.4f&longitude=%.4f&daily=temperature_2m_max,temperature_2m_min&start_date=%s&end_date=%s&temperature_unit=fahrenheit",
		forecastBaseURL,
		lat, lon,
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)

	return fetchDailyTemps(u)
}

// fetchHistorical uses the Open-Meteo Historical Archive API for past dates
// or for "typical" temperatures (same dates from a prior year).
func fetchHistorical(lat, lon float64, start, end time.Time) (minTemp, maxTemp float64, err error) {
	u := fmt.Sprintf(
		"%s/v1/archive?latitude=%.4f&longitude=%.4f&daily=temperature_2m_max,temperature_2m_min&start_date=%s&end_date=%s&temperature_unit=fahrenheit",
		archiveBaseURL,
		lat, lon,
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)

	return fetchDailyTemps(u)
}

// fetchDailyTemps fetches daily temperature data from the given URL and
// returns the overall minimum low and maximum high across all days.
func fetchDailyTemps(u string) (minTemp, maxTemp float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("weather request failed: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("weather request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var weather dailyWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return 0, 0, fmt.Errorf("failed to decode weather response: %w", err)
	}

	if len(weather.Daily.TemperatureMin) == 0 || len(weather.Daily.TemperatureMax) == 0 {
		return 0, 0, fmt.Errorf("no temperature data returned")
	}

	minTemp = weather.Daily.TemperatureMin[0]
	maxTemp = weather.Daily.TemperatureMax[0]

	for _, t := range weather.Daily.TemperatureMin {
		if t < minTemp {
			minTemp = t
		}
	}
	for _, t := range weather.Daily.TemperatureMax {
		if t > maxTemp {
			maxTemp = t
		}
	}

	return minTemp, maxTemp, nil
}
