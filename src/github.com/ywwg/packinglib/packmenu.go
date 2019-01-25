package packinglib

// PackMenuType is a type for items in a promptui selection.
type PackMenuType int

const (
	// MenuCategory is the type for a packing category
	MenuCategory PackMenuType = iota

	// MenuPackable is the type for individual packable items
	MenuPackable

	// MenuProperty is the type for packing properties.
	MenuProperty

	// MenuAction is the type for prompt interaction operations (back, quit, etc)
	MenuAction
)

// PackMenuItem is a struct for each item in a promptui selection
type PackMenuItem struct {
	// Name is the display name with formatting whitespace
	Name string

	// Type determines how the item will be handled
	Type PackMenuType

	// Code is the machine-readable code for determining how the item will be
	// handled.
	Code string
}

// NewMenuItem creates a new menu item with the given name, type, and code.
func NewMenuItem(name string, t PackMenuType, code string) PackMenuItem {
	displayName := name
	if t == MenuPackable {
		displayName = " " + name
	}

	return PackMenuItem{
		Name: displayName,
		Type: t,
		Code: code,
	}
}

// Equals return true if this menu item equals the other.
func (i PackMenuItem) Equals(other PackMenuItem) bool {
	return i.Type == other.Type && i.Code == other.Code
}
