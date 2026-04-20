package packinglib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructRegistryClone(t *testing.T) {
	src := NewTestRegistry()
	clone := src.Clone()

	require.NotSame(t, src, clone, "clone must be a distinct registry instance")
	require.ElementsMatch(t, src.ContextList(), clone.ContextList())

	// Same property keys
	srcProps := make([]Property, 0, len(src.AllProperties()))
	for p := range src.AllProperties() {
		srcProps = append(srcProps, p)
	}
	cloneProps := make([]Property, 0, len(clone.AllProperties()))
	for p := range clone.AllProperties() {
		cloneProps = append(cloneProps, p)
	}
	require.ElementsMatch(t, srcProps, cloneProps)

	// Same item names per category
	for cat, items := range src.AllItems() {
		srcNames := make([]string, 0, len(items))
		for _, it := range items {
			srcNames = append(srcNames, it.Name())
		}
		cloneItems, ok := clone.AllItems()[cat]
		require.True(t, ok, "clone missing category %s", cat)
		cloneNames := make([]string, 0, len(cloneItems))
		for _, it := range cloneItems {
			cloneNames = append(cloneNames, it.Name())
		}
		require.ElementsMatch(t, srcNames, cloneNames, "category %s item names", cat)
	}
}

func TestStructRegistryCloneIndependentMaps(t *testing.T) {
	src := NewTestRegistry()
	clone := src.Clone()

	clone.RegisterProperty("new-prop", "only in clone")
	require.False(t, src.HasProperty("new-prop"), "source must not see property added on clone")
	require.True(t, clone.HasProperty("new-prop"))

	clone.RegisterItems("clone-only-category", []*Item{NewItem("clone-only-item", nil, nil)})
	_, srcHas := src.AllItems()["clone-only-category"]
	require.False(t, srcHas, "source must not see category added on clone")
}

func TestStructRegistryCloneIndependentItems(t *testing.T) {
	src := NewTestRegistry()
	clone := src.Clone()

	// Pack the first item in the clone's "clothing" category.
	cloneClothing := clone.AllItems()["clothing"]
	require.NotEmpty(t, cloneClothing)
	cloneClothing[0].Pack(true)

	// Find the same-named item in the source and verify it is still unpacked.
	name := cloneClothing[0].Name()
	srcClothing := src.AllItems()["clothing"]
	var found *Item
	for _, it := range srcClothing {
		if it.Name() == name {
			found = it
			break
		}
	}
	require.NotNil(t, found, "source clothing must still contain %s", name)
	require.False(t, found.Packed(), "packing on clone must not affect source item %s", name)
}

func TestStructRegistryCloneDupeChecker(t *testing.T) {
	src := NewTestRegistry()
	clone := src.Clone()

	// Register a new unique item in the clone.
	clone.RegisterProperty("isolate", "isolate test")
	require.NotPanics(t, func() {
		clone.RegisterItems("extra", []*Item{NewItem("isolate-probe", []string{"isolate"}, nil)})
	})

	// The source must still allow registering THAT SAME name without panic,
	// proving dupeChecker was not shared with the clone.
	src.RegisterProperty("isolate", "isolate test")
	require.NotPanics(t, func() {
		src.RegisterItems("extra", []*Item{NewItem("isolate-probe", []string{"isolate"}, nil)})
	})
}
