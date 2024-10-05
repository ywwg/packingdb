package packinglib

import "fmt"

const (
	fmt_version = 1
)

type YamlTrip struct {
	Version    int         `yaml:"version"`
	Name       string      `yaml:"name"`
	Nights     int         `yaml:"nights"`
	TempMin    int         `yaml:"temp_min"`
	TempMax    int         `yaml:"temp_max"`
	Properties []string    `yaml:"properties,omitempty"`
	PackList   []*YamlItem `yaml:"pack_list,omitempty"`
}

// // YamlPackItem is a minimal verison of YamlItem that only stores what
// // is needed to log if it is packed.
// type YamlPackItem struct {
// }

func FromTrip(t *Trip) *YamlTrip {
	yt := &YamlTrip{
		Version: fmt_version,
		Name:    t.contextName,
		Nights:  t.C.Nights,
		TempMin: t.C.TemperatureMin,
		TempMax: t.C.TemperatureMax,
	}

	for p := range t.C.Properties {
		yt.Properties = append(yt.Properties, string(p))
	}

	for _, il := range t.packList {
		for _, i := range il {
			yt.PackList = append(yt.PackList, PackItem(i))
		}
	}

	return yt
}

func (yt *YamlTrip) AsTrip(r Registry) (*Trip, error) {
	c, err := NewContext(r, yt.Name, yt.Nights, yt.TempMin, yt.TempMax, yt.Properties)
	if err != nil {
		return nil, err
	}

	t, err := NewTripFromCustomContext(r, c)
	if err != nil {
		return nil, err
	}

	for _, i := range yt.PackList {
		t.Pack(i.Name, *i.Packed)
	}

	return t, nil
}

type YamlItemList struct {
	Name  string     `yaml:"name"`
	Items []YamlItem `yaml:"items"`
}

func (yl *YamlItemList) AsItemList() (*ItemList, error) {
	il := &ItemList{
		Name: yl.Name,
	}

	for _, yi := range yl.Items {
		item := &Item{
			name:  yi.Name,
			units: yi.Units,
		}

		item.prerequisites = buildPropertySet(yi.Allow, yi.Disallow)

		if yi.TempMax != nil {
			if yi.TempMin != nil {
				item.TemperatureRange(*yi.TempMin, *yi.TempMax)
			} else {
				return nil, fmt.Errorf("both temp_max and temp_min must be set if one is")
			}
		} else if yi.TempMin != nil {
			return nil, fmt.Errorf("both temp_max and temp_min must be set if one is")
		}

		if yi.PerDay != nil {
			item.Consumable(*yi.PerDay)
		}

		if yi.Max != nil {
			item.Max(*yi.Max)
		}

		if yi.CustomFuncName != "" {
			// TODO switch goes here
			return nil, fmt.Errorf("unknown custom func %s", yi.CustomFuncName)
		}

		il.Items = append(il.Items, item)
	}

	return il, nil
}

type YamlItem struct {
	Name     string   `yaml:"name"`
	Units    string   `yaml:"units,omitempty"`
	Allow    []string `yaml:"allow,omitempty"`
	Disallow []string `yaml:"disallowed,omitempty"`

	TempMin        *int     `yaml:"temp_min,omitempty"`
	TempMax        *int     `yaml:"temp_max,omitempty"`
	PerDay         *float64 `yaml:"perday,omitempty"`
	Max            *float64 `yaml:"max,omitempty"`
	CustomFuncName string   `yaml:"custom_func_name,omitempty"`

	Packed *bool `yaml:"packed,omitempty"`
}

// PackItem returns an extremely minimal version of YamlItem that only
// records the name, and whether the item is packed.
func PackItem(i *Item) *YamlItem {
	return &YamlItem{
		Name:   i.name,
		Packed: &i.packed,
	}
}
