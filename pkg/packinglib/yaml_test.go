package packinglib

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestYamlItemList(t *testing.T) {
	yil := &YamlItemList{
		Name: "listname",
		Items: []YamlItem{
			{
				Name: "item1",
			},
			{
				Name:           "item2",
				Units:          "myunits",
				Allow:          []string{"allow1", "allow2"},
				Disallow:       []string{"disallow1", "disallow2"},
				TempMin:        Int(0),
				TempMax:        Int(100),
				PerDay:         proto.Float64(1.5),
				Max:            proto.Float64(10.2),
				CustomFuncName: "mycustomfunc",
				Packed:         proto.Bool(false),
			},
			{
				Name:     "item3",
				Units:    "myunits2",
				Allow:    []string{"allow3", "allow4"},
				Disallow: []string{"disallow5", "disallow6"},
				TempMin:  Int(0),
				TempMax:  Int(100),
				PerDay:   proto.Float64(1.5),
				Max:      proto.Float64(10.2),
				Packed:   proto.Bool(true),
			},
		},
	}

	il, err := yil.AsItemList()
	require.NoError(t, err)

	require.Equal(t, ItemList{
		Name: "listname",
		Items: []*Item{
			{
				name:  "item2",
				units: "myunits",
				prerequisites: PropertySet{
					"allow1":    true,
					"allow2":    true,
					"disallow1": false,
					"disallow2": false,
				},
				// XXXX mutators...
			},
		},
	}, il)
}

func Int(v int) *int { return &v }
