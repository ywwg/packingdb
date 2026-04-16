# External Integrations

**Analysis Date:** 2026-04-16

## APIs & External Services

**Weather & Location:**
- Open-Meteo Geocoding API - Location name to coordinates resolution
  - SDK/Client: Native Go `net/http` with custom timeout (10 seconds)
  - Base URL: `https://geocoding-api.open-meteo.com/v1/search`
  - Used in: `pkg/weather/weather.go` - `geocode()`, `GeocodeSuggestions()` functions
  - Public API - No auth required
  - Features: Location autocomplete filtering by state/province or country

- Open-Meteo Forecast API - Weather forecast (up to 16 days ahead)
  - SDK/Client: Native Go `net/http`
  - Base URL: `https://api.open-meteo.com/v1/forecast`
  - Used in: `pkg/weather/weather.go` - `fetchForecast()` function
  - Public API - No auth required
  - Response includes: Daily temperature min/max in Fahrenheit

- Open-Meteo Historical Archive API - Past/typical temperatures
  - SDK/Client: Native Go `net/http`
  - Base URL: `https://archive-api.open-meteo.com/v1/archive`
  - Used in: `pkg/weather/weather.go` - `fetchHistorical()` function
  - Public API - No auth required
  - Purpose: Fetch "typical" temperatures from same dates in previous year for future trips beyond 16-day forecast window

**Weather Logic:**
- Forecast API selected for trips within 16 days ahead
- Historical API (previous year same dates) selected for trips beyond 16 days
- All calls time-bounded to 10 seconds (configurable in `pkg/weather/weather.go` line 17)

## Data Storage

**Databases:**
- None - No database service

**File Storage:**
- Local filesystem only
  - Format: YAML
  - Storage location: `./public/trips/` (configurable via `-trips` flag)
  - File naming: Sanitized trip names with auto-collision handling via suffix (`sanitizeFilename()`)
  - Trip structure: `pkg/packinglib/yaml.go` defines `YamlTrip` format
  - Persistence: Background goroutine saves trips periodically in `pkg/server/persist.go`

**Caching:**
- In-memory only
  - Trip objects cached in `pkg/server/server.go` - `trips map[string]*Trip`
  - Name-to-filename mapping cached for performance

## Authentication & Identity

**Auth Provider:**
- None - No authentication
- All APIs and endpoints are unauthenticated
- Public Web interface - No access control

## Monitoring & Observability

**Error Tracking:**
- None

**Logs:**
- Custom structured logging via `pkg/logger/` package
  - Logging levels: debug, info, warn, error (configurable at startup)
  - Called throughout: server initialization, trip creation, persistence, API errors
  - Implementation: `pkg/logger/` (usage seen in `cmd/packingweb/main.go`, `pkg/server/`)

## CI/CD & Deployment

**Hosting:**
- Self-hosted only (no cloud platform integration)
- Single Go binary execution model

**CI Pipeline:**
- None detected - No CI/CD configuration files present

**Deployment Model:**
- Manual: `go run ./cmd/packingweb` or build binary and execute
- Data persists in local `./public/trips/` directory
- Graceful shutdown handling: SIGINT and SIGTERM signals (see `cmd/packingweb/main.go`)

## Environment Configuration

**Required env vars:**
- None - All configuration via command-line flags or defaults

**Secrets location:**
- Not applicable - No secrets management needed
- All APIs (Open-Meteo) are public and require no keys

## Webhooks & Callbacks

**Incoming:**
- None

**Outgoing:**
- None

## CORS & Cross-Origin

**CORS Headers:**
- Enabled in `pkg/server/routes.go` - `corsMiddleware()` function
- Allows methods: GET, POST, PUT, DELETE, OPTIONS
- Allows all origins (no origin restrictions)
- Purpose: Enable frontend calls from any origin

---

*Integration audit: 2026-04-16*
