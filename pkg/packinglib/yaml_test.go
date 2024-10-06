package packinglib

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestYamlItemList(t *testing.T) {
	buf := []byte(`name: listname
items:
    - name: item1
    - name: item2
      units: myunits
      allow:
        - allow1
        - allow2
      disallowed:
        - disallow1
        - disallow2
      temp_min: 0
      temp_max: 100
      perday: 1.5
      max: 10.2
      packed: false
    - name: item3
      units: myunits2
      allow:
        - allow3
        - allow4
      disallowed:
        - disallow5
        - disallow6
      temp_min: 0
      temp_max: 100
      perday: 1.5
      packed: true`)
	reader := bytes.NewBuffer(buf)
	dec := yaml.NewDecoder(reader)
	var yil *YamlItemList
	require.NoError(t, dec.Decode(&yil))

	il, err := yil.AsItemList()
	require.NoError(t, err)

	expectedAttrs := []struct {
		Name          string
		Units         string
		Prerequisites PropertySet
		Mutators      []string
	}{
		{
			Name:          "item1",
			Prerequisites: PropertySet{},
		},
		{
			Name:  "item2",
			Units: "myunits",
			Prerequisites: PropertySet{
				"allow1":    true,
				"allow2":    true,
				"disallow1": false,
				"disallow2": false,
			},
			Mutators: []string{"temperatureMutator", "consumableMutator", "maxCountMutator"},
		},
		{
			Name:  "item3",
			Units: "myunits2",
			Prerequisites: PropertySet{
				"allow3":    true,
				"allow4":    true,
				"disallow5": false,
				"disallow6": false,
			},
			Mutators: []string{"temperatureMutator", "consumableMutator"},
		},
	}

	require.Equal(t, len(expectedAttrs), len(il.Items))

	for i, item := range il.Items {
		require.Equalf(t, expectedAttrs[i].Name, item.Name(), "element %d mismatch", i)
		require.Equalf(t, expectedAttrs[i].Units, item.units, "element %d mismatch", i)
		require.Equalf(t, expectedAttrs[i].Prerequisites, item.Prerequisites(), "element %d mismatch", i)
		require.Equalf(t, len(expectedAttrs[i].Mutators), len(item.mutators), "mutator len mismatch in element %d", i)
		for j, m := range item.mutators {
			require.Equal(t, expectedAttrs[i].Mutators[j], m.Name())
		}
	}
}

func TestYamlTrip(t *testing.T) {
	buf := []byte(`version: 1
name: berlin work trip
nights: 5
temp_min: 46
temp_max: 66
properties:
    - berlin work trip
    - Business
    - Flight
    - International
pack_list:
    - name: work laptop
      packed: true
    - name: socks
      packed: true
`)
	reader := bytes.NewBuffer(buf)
	dec := yaml.NewDecoder(reader)
	var yt *YamlTrip
	require.NoError(t, dec.Decode(&yt))

	var r Registry = NewStructRegistry()
	PopulateRegistry(r)

	got, err := yt.AsTrip(r)
	require.NoError(t, err)

	require.Equal(t, "berlin work trip", got.C.Name)
	require.Equal(t, 5, got.C.Nights)
	require.Equal(t, 46, got.C.TemperatureMin)
	require.Equal(t, 66, got.C.TemperatureMax)
	require.Equal(t, PropertySet{
		"Business":         true,
		"Flight":           true,
		"International":    true,
		"berlin work trip": true,
	}, got.C.Properties)

	require.Equal(t, "berlin work trip", got.contextName)

	require.Equal(t, 1, len(got.packList))
	require.Equal(t, "work laptop", got.packList["business"][0].Name())
	require.Equal(t, float64(1), got.packList["business"][0].count)
	require.True(t, got.packList["business"][0].packed)

	require.Equal(t, "socks", got.packList["business"][1].Name())
	require.Equal(t, float64(5), got.packList["business"][1].count)
	require.True(t, got.packList["business"][1].packed)

	// Finally test the round trip.
	var buf2 []byte
	writer := bytes.NewBuffer(buf2)
	yt = FromTrip(got)

	enc := yaml.NewEncoder(writer)
	require.NoError(t, enc.Encode(yt))

	require.Equal(t, string(buf), writer.String())
}

func PopulateRegistry(r Registry) {
	r.RegisterProperty("Business", "business cat")
	r.RegisterProperty("Flight", "wheee")
	r.RegisterProperty("International", "le whee")
	r.RegisterProperty("camping", "tent")

	r.RegisterItems("business", []*Item{
		NewItem("work laptop", []string{"Business"}, nil),
		NewItem("socks", []string{"Flight"}, []string{"camping"}).Consumable(1),
	})
}

func Int(v int) *int { return &v }
