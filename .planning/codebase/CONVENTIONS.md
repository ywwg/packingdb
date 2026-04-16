# Coding Conventions

**Analysis Date:** 2026-04-16

## Naming Patterns

**Files:**
- Package main files follow pattern `packingdb.go`, `packingcli.go`, `main.go`
- Test files suffixed with `_test.go` (`item_test.go`, `yaml_test.go`, `weather_test.go`)
- Public API functions and types use PascalCase: `NewItem`, `NewContext`, `Result`, `Registry`

**Functions:**
- Public functions use PascalCase: `NewItem()`, `Satisfies()`, `AdjustCount()`, `String()`
- Private functions and methods use camelCase: `geocode()`, `fetchForecast()`, `addProperty()`
- Getter methods use simple PascalCase: `Name()`, `Count()`, `Packed()`, `Prerequisites()`

**Variables:**
- Local variables and parameters use camelCase: `startDate`, `endDate`, `tripNames`, `tmpDir`
- Package-level state variables use camelCase: `httpClient`, `dupeChecker`
- Constants mixed case based on visibility: `yaml_version` (private), `TemperatureMin` (public struct field)

**Types:**
- Interfaces end with "er" or descriptive names: `packMutator`, `Registry`
- Struct types use PascalCase: `Item`, `Context`, `Property`, `PropertySet`, `Result`
- Specialized types (often based on primitives) use PascalCase: `Category` (type string), `Property` (type string)

## Code Style

**Formatting:**
- Standard Go formatting (follows `gofmt` conventions)
- No explicit linter configuration file found; assume gofmt defaults
- Lines appear to follow standard Go width conventions

**Linting:**
- No `.golangci.yml` or `.eslintrc` found
- Code uses standard Go idioms (error handling, interface design)
- No strict linting enforcement detected

## Import Organization

**Order:**
1. Standard library imports (stdlib): `fmt`, `errors`, `sort`, `os`, `log`, etc.
2. External packages: `github.com/manifoldco/promptui`, `gopkg.in/yaml.v3`, etc.
3. Local project imports: `github.com/ywwg/packingdb/pkg/...`

**Path Aliases:**
- No import aliases used in codebase
- Full paths: `"github.com/ywwg/packingdb/pkg/packinglib"`

**Example from `routes.go`:**
```go
import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/ywwg/packingdb/pkg/logger"
	"github.com/ywwg/packingdb/pkg/packinglib"
)
```

## Error Handling

**Patterns:**
- Explicit error checking: `if err != nil { return nil, err }`
- Error wrapping with context: `fmt.Errorf("geocoding failed: %w", err)` (uses `%w` for error chain)
- Validation errors with descriptive messages: `fmt.Errorf("didn't find property, is it registered?: %s", prop)`
- Function returns `(value, error)` tuple: `Lookup(...) (*Result, error)`

**Strategy in `weather.go`:**
- Each function checks errors immediately
- Errors are wrapped with human-readable context (line 82: `"geocoding failed: %w"`)
- No silent failures or error suppression

**Strategy in `registry.go`:**
- Panics for invariant violations (unrecoverable errors): `panic(fmt.Sprintf("Duplicate context: %s", c.Name))`
- Returns errors for recoverable cases: `fmt.Errorf("unknown context: %s", name)`

## Logging

**Framework:** 
- Standard `log` package for CLI: `log.Fatal()`, `log.Print()`
- Custom package `pkg/logger` for server (imported in `routes.go`): `logger.Info("Created new trip", "name", req.Name, "file", filename)`
- HTTP handlers via `server.respondError(w, message, status)`

**Patterns:**
- CLI tools use `log.Fatal()` for fatal conditions
- Server logs key operations with structured key-value pairs: `logger.Info(msg, key, value, ...)`
- No error logging visible in public code (likely handled at call site)

## Comments

**When to Comment:**
- Top-of-file package documentation: `// Package weather provides temperature lookup...` (line 1 in weather.go)
- Complex logic with non-obvious intent: temperatureMutator behavior, property satisfaction rules
- Public API documentation above exported types and functions

**JSDoc/TSDoc:**
- Not used (Go project, uses GoDoc conventions)
- Go-style comments above exported functions: `// Lookup geocodes the location...`
- Struct field comments: `// Name of the item.` (in item.go lines 12-23)

**Examples from codebase:**
```go
// Item
type Item struct {
	// Name of the item.
	name string
	
	// count is the number of this thing that should get packed.
	count float64
	
	// Prerequisites is a set of all properties that the context must have for this item to appear.
	prerequisites PropertySet
}
```

```go
// Lookup geocodes the location and fetches weather data for the given date
// range. Dates must be today or in the future. If the dates are within the
// 16-day forecast window, it uses the forecast API. If the dates are further
// out, it looks up historical data from the same dates in the previous year
// to get "typical" temperatures.
func Lookup(location string, startDate, endDate time.Time) (*Result, error) {
```

## Function Design

**Size:**
- Small, focused functions preferred
- `Lookup()` is 57 lines (longest in weather.go) but has clear sections
- Most mutator methods are 3-10 lines
- Clear separation of concerns: `fetchForecast()`, `fetchHistorical()`, `fetchDailyTemps()`

**Parameters:**
- Functions accept concrete types, not interfaces (except `Registry` where needed)
- Multiple return values: `(value, error)` pattern
- Named return types not used; callers don't rely on them

**Return Values:**
- Always include error as second return: `func Lookup(...) (*Result, error)`
- Pointer receivers for mutation: `func (i *Item) Units(u string) *Item { ... return i }`
- Fluent builder pattern: methods return `*Item` to enable chaining

**Example fluent pattern from `item.go`:**
```go
func (i *Item) Units(u string) *Item {
	i.units = u
	return i
}

func (i *Item) TemperatureRange(tMin, tMax int) *Item {
	i.mutators = append(i.mutators, &temperatureMutator{tMin, tMax})
	return i
}

// Called as:
i := NewItem("mytempitem", []string{"prop1"}, []string{"prop3"})
i.TemperatureRange(0, 100)
i.Consumable(2)
i.Max(5)
```

## Module Design

**Exports:**
- Public types are UPPERCASE: `Item`, `Context`, `Registry`, `Result`
- Public functions are UPPERCASE: `NewItem()`, `Lookup()`, `NewContext()`
- Private types and functions are lowercase: `packMutator`, `geocode()`, `temperatureMutator`

**Barrel Files:**
- No barrel files or re-exports; packages export directly
- `pkg/packinglib/` exports all its types directly
- `pkg/items/` contains item registration but no re-exports

**File Organization:**
- Related types grouped by concern: `item.go` contains all Item-related logic + mutators
- Test files in same package: `item_test.go`, `yaml_test.go`, `weather_test.go`
- No separate internal/ packages; relying on visibility conventions

## Typical Code Patterns

**Builder Pattern:**
```go
// From item_test.go
i := NewItem("myitem", []string{"prop1"}, []string{"prop3"})
i.Units("pounds")
i.Pack(true)
i.AdjustCount(&basicContext)
```

**Mutator Interface Implementation:**
```go
type packMutator interface {
	Name() string
	AdjustCount(c *Context, count float64) float64
	Priority() int
}

// Then concrete implementations with pointer receivers:
func (m *temperatureMutator) AdjustCount(c *Context, count float64) float64 {
	// ...
}
```

**Handler Pattern (from routes.go):**
```go
func (s *Server) listTripsHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	// ... read data
	s.mu.RUnlock()
	s.respondJSON(w, responseData)
}
```

---

*Convention analysis: 2026-04-16*
