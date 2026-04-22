<!-- GSD:project-start source:PROJECT.md -->
## Project

**PackingDB: Multi-Trip State Isolation**

A packing list manager with CLI and web interfaces. Users create trips with properties (camping, cold weather, etc.), and the system generates a filtered packing list. The web server (`packingweb`) serves a REST API and static frontend. This milestone fixes a state bleed bug and adds proper multi-trip support to the backend.

**Core Value:** Each trip's packing state must be completely independent. Packing items in Trip A must never affect Trip B.

### Constraints

- **Backend only**: No frontend changes. API contract stays the same.
- **Existing persistence**: Dirty trip sync already writes to disk every 30s. Leverage it, don't replace it.
- **Brownfield**: This is an existing Go codebase. Follow existing patterns (chi router, packinglib structure, YAML/CSV serialization).
<!-- GSD:project-end -->

<!-- GSD:stack-start source:codebase/STACK.md -->
## Technology Stack

## Languages
- Go 1.25.7 - Backend server and CLI applications (`cmd/packingweb`, `cmd/packingcli`, `cmd/packingdb`)
- JavaScript (vanilla) - Frontend logic (`static/app.js`)
- HTML5 - Web interface (`static/index.html`)
- CSS3 - Styling via Tailwind CDN (`static/styles.css`)
- YAML - Trip data serialization (`pkg/packinglib/yaml.go`)
## Runtime
- Go 1.25.7 runtime
- Go modules (go.mod)
- Lockfile: `go.sum` present
## Frameworks
- chi/v5 v5.2.4 - HTTP router and middleware (`github.com/go-chi/chi/v5`)
- Alpine.js 3.15.8 - Lightweight reactive frontend framework
- Tailwind CSS (dynamic CDN build) - Utility-first CSS framework
- testify v1.9.0 - Assertion library and test utilities (`github.com/stretchr/testify`)
- Go built-in `testing` package
- promptui v0.9.0 - Interactive terminal UI for CLI prompts (`github.com/manifoldco/promptui`)
- Go build tools (no dedicated build system)
## Key Dependencies
- chi/v5 v5.2.4 - HTTP routing, fundamental to REST API server
- promptui v0.9.0 - Interactive CLI experience for traditional interface
- golang.org/x/crypto v0.22.0 - Cryptographic utilities (in transitive chain)
- gopkg.in/yaml.v3 v3.0.1 - YAML parsing/serialization for trip data storage
- golang.org/x/sys v0.19.0 - System-level utilities (indirect, via promptui)
- golang.org/x/term v0.19.0 - Terminal utilities (indirect, via promptui)
## Configuration
- Configured via command-line flags in `cmd/packingweb/main.go`:
- `go.mod` and `go.sum` files
- No build configuration files (Makefile, gradle, etc.)
- Shell scripts provided: `packingdb.sh`, `packingcli.sh` for quick execution
## Platform Requirements
- Go 1.25.7 installed
- GOPATH configuration (noted in README.md as required for module resolution)
- Standard Unix build tools (make, etc.) implicit in shell scripts
- Go 1.25.7 runtime (single binary deployment)
- Directories created at runtime: `./public/trips` (trip data), static file serving from `./static`
- HTTP port exposed (configurable, default 8080)
- Local filesystem access for trip file persistence
- Modern browser with ES6 support
- JavaScript enabled
- No package manager (CDN-delivered dependencies)
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

## Naming Patterns
- Package main files follow pattern `packingdb.go`, `packingcli.go`, `main.go`
- Test files suffixed with `_test.go` (`item_test.go`, `yaml_test.go`, `weather_test.go`)
- Public API functions and types use PascalCase: `NewItem`, `NewContext`, `Result`, `Registry`
- Public functions use PascalCase: `NewItem()`, `Satisfies()`, `AdjustCount()`, `String()`
- Private functions and methods use camelCase: `geocode()`, `fetchForecast()`, `addProperty()`
- Getter methods use simple PascalCase: `Name()`, `Count()`, `Packed()`, `Prerequisites()`
- Local variables and parameters use camelCase: `startDate`, `endDate`, `tripNames`, `tmpDir`
- Package-level state variables use camelCase: `httpClient`, `dupeChecker`
- Constants mixed case based on visibility: `yaml_version` (private), `TemperatureMin` (public struct field)
- Interfaces end with "er" or descriptive names: `packMutator`, `Registry`
- Struct types use PascalCase: `Item`, `Context`, `Property`, `PropertySet`, `Result`
- Specialized types (often based on primitives) use PascalCase: `Category` (type string), `Property` (type string)
## Code Style
- Standard Go formatting (follows `gofmt` conventions)
- No explicit linter configuration file found; assume gofmt defaults
- Lines appear to follow standard Go width conventions
- No `.golangci.yml` or `.eslintrc` found
- Code uses standard Go idioms (error handling, interface design)
- No strict linting enforcement detected
## Import Organization
- No import aliases used in codebase
- Full paths: `"github.com/ywwg/packingdb/pkg/packinglib"`
## Error Handling
- Explicit error checking: `if err != nil { return nil, err }`
- Error wrapping with context: `fmt.Errorf("geocoding failed: %w", err)` (uses `%w` for error chain)
- Validation errors with descriptive messages: `fmt.Errorf("didn't find property, is it registered?: %s", prop)`
- Function returns `(value, error)` tuple: `Lookup(...) (*Result, error)`
- Each function checks errors immediately
- Errors are wrapped with human-readable context (line 82: `"geocoding failed: %w"`)
- No silent failures or error suppression
- Panics for invariant violations (unrecoverable errors): `panic(fmt.Sprintf("Duplicate context: %s", c.Name))`
- Returns errors for recoverable cases: `fmt.Errorf("unknown context: %s", name)`
## Logging
- Standard `log` package for CLI: `log.Fatal()`, `log.Print()`
- Custom package `pkg/logger` for server (imported in `routes.go`): `logger.Info("Created new trip", "name", req.Name, "file", filename)`
- HTTP handlers via `server.respondError(w, message, status)`
- CLI tools use `log.Fatal()` for fatal conditions
- Server logs key operations with structured key-value pairs: `logger.Info(msg, key, value, ...)`
- No error logging visible in public code (likely handled at call site)
## Comments
- Top-of-file package documentation: `// Package weather provides temperature lookup...` (line 1 in weather.go)
- Complex logic with non-obvious intent: temperatureMutator behavior, property satisfaction rules
- Public API documentation above exported types and functions
- Not used (Go project, uses GoDoc conventions)
- Go-style comments above exported functions: `// Lookup geocodes the location...`
- Struct field comments: `// Name of the item.` (in item.go lines 12-23)
## Function Design
- Small, focused functions preferred
- `Lookup()` is 57 lines (longest in weather.go) but has clear sections
- Most mutator methods are 3-10 lines
- Clear separation of concerns: `fetchForecast()`, `fetchHistorical()`, `fetchDailyTemps()`
- Functions accept concrete types, not interfaces (except `Registry` where needed)
- Multiple return values: `(value, error)` pattern
- Named return types not used; callers don't rely on them
- Always include error as second return: `func Lookup(...) (*Result, error)`
- Pointer receivers for mutation: `func (i *Item) Units(u string) *Item { ... return i }`
- Fluent builder pattern: methods return `*Item` to enable chaining
## Module Design
- Public types are UPPERCASE: `Item`, `Context`, `Registry`, `Result`
- Public functions are UPPERCASE: `NewItem()`, `Lookup()`, `NewContext()`
- Private types and functions are lowercase: `packMutator`, `geocode()`, `temperatureMutator`
- No barrel files or re-exports; packages export directly
- `pkg/packinglib/` exports all its types directly
- `pkg/items/` contains item registration but no re-exports
- Related types grouped by concern: `item.go` contains all Item-related logic + mutators
- Test files in same package: `item_test.go`, `yaml_test.go`, `weather_test.go`
- No separate internal/ packages; relying on visibility conventions
## Typical Code Patterns
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

## Pattern Overview
- Domain-driven design with core packing logic isolated from transport layers
- Multiple entry points (CLI, TUI, HTTP API) all sharing the same core library
- Registry pattern for managing extensible item and property definitions
- Background persistence for server-side changes
- File-based storage with CSV and YAML support
## Layers
- Purpose: Pure business logic for managing trips, items, properties, and packing lists
- Location: `pkg/packinglib/`
- Contains: Trip orchestration, item filtering, context management, property sets, file loading/saving
- Depends on: Standard library only (golang.org/x/crypto for password hashing in some cases, gopkg.in/yaml.v3 for serialization)
- Used by: All other packages and entry points
- Purpose: Defines and registers available items, properties, and pre-built contexts
- Location: `pkg/contexts/contexts.go`
- Contains: Item definitions from `pkg/items/*`, property list definitions, context templates (Firefly, Cape, Tiny House variants)
- Depends on: `pkg/packinglib`, `pkg/items`
- Used by: Server startup and CLI initialization
- Purpose: Categorized item definitions with property prerequisites
- Location: `pkg/items/` (split by category: clothing.go, electrical.go, toiletries.go, etc.)
- Contains: Item lists with conditional logic based on trip properties
- Depends on: `pkg/packinglib` (Item, Category types)
- Used by: Registry during initialization
- Purpose: REST API exposure of domain operations with request/response marshaling
- Location: `pkg/server/`
- Contains: HTTP handlers, routing, JSON serialization, file-to-trip mapping
- Depends on: `pkg/packinglib`, `pkg/logger`, standard http library, go-chi router
- Used by: HTTP server entry point (`cmd/packingweb`)
- Purpose: Interactive terminal-based interface for trip management
- Location: `cmd/packingdb/packingdb.go`, `cmd/packingcli/packingcli.go`
- Contains: Prompt UI menus, console I/O, trip orchestration via promptui
- Depends on: `pkg/packinglib`, `pkg/contexts`, manifoldco/promptui
- Used by: Terminal users
- Purpose: HTTP server entry point that initializes registry and starts API
- Location: `cmd/packingweb/main.go`
- Contains: Flag parsing, logger initialization, server creation, graceful shutdown handling
- Depends on: `pkg/server`, `pkg/packinglib`, `pkg/contexts`, `pkg/logger`
- Used by: Users running the web interface
- Logger: `pkg/logger/logger.go` - Structured logging with level support
- Weather: `pkg/weather/weather.go` - Temperature/weather lookups (currently supports location-based weather API calls)
## Data Flow
- In-memory: Trips cached in Server.trips map (keyed by trip name)
- On-disk: Trips saved to `public/trips/` directory as .yml or .csv files
- Dirty tracking: dirtyTrips map holds names of trips needing persistence
- Background sync: Persistence goroutine runs every 30s, writes only dirty trips
## Key Abstractions
- Purpose: Represents a complete packing scenario with context, items, and pack state
- Examples: `cmd/packingweb/main.go` line 37, `pkg/server/routes.go` line 81
- Pattern: Mutable value object that holds filtered ItemList and applies mutations (Pack, AddProperty)
- Purpose: Describes conditions that determine which items are relevant (nights, temperature, properties)
- Examples: Pre-built contexts in `pkg/contexts/contexts.go` (fireflyContext, capeContext, tinyhouseSpring, etc.)
- Pattern: Immutable by design, holds PropertySet that drives item filtering
- Purpose: Service locator for all items, properties, and context definitions
- Examples: `pkg/packinglib/registry.go` interface definition, `pkg/contexts/contexts.go` PopulateRegistry()
- Pattern: Interface (Registry) with in-memory implementation (StructRegistry), populated at startup
- Purpose: Individual thing to pack with conditional inclusion based on context
- Examples: `pkg/items/clothing.go`, `pkg/items/toiletries.go` item definitions
- Pattern: Immutable item name/description with mutable packed state and prerequisite matching
- Purpose: Efficient set-like structure for property matching in Item.Satisfies()
- Pattern: True=allow, False=disallow; enables flexible opt-in and opt-out filtering
## Entry Points
- Location: `cmd/packingweb/main.go`
- Triggers: User runs binary; listens on :8080
- Responsibilities: HTTP server lifecycle, static file serving, API routing, background persistence
- Location: `cmd/packingdb/packingdb.go`
- Triggers: User runs `./packingdb <tripfile.yml>`
- Responsibilities: Interactive menus via promptui, trip manipulation, file I/O
- Location: `cmd/packingcli/packingcli.go`
- Triggers: User runs `packingcli` (older interface)
- Responsibilities: Non-interactive command-line trip creation
## Error Handling
- File operations return errors (LoadFromFile, SaveToFile)
- Trip creation returns errors (NewTrip, NewContext)
- Property/Item lookup returns errors (AddProperty, Pack)
- HTTP handlers use respondError() with appropriate status codes (400 BadRequest, 409 Conflict, 500 InternalServerError)
- Registry lookups return errors (GetContext, HasProperty)
- Logger.Fatal() used for fatal startup conditions (server creation, directory creation)
## Cross-Cutting Concerns
- Framework: `pkg/logger/logger.go` provides Debug, Info, Warn, Error, Fatal functions
- Configured at startup via `-log-level` flag (default: info)
- Used for: Server startup, trip persistence, scan results, API requests
- Request body validation in HTTP handlers (name required, nights > 0)
- Trip file parsing validates YAML/CSV format and version
- Property existence checked against registry (HasProperty)
- Temperature range constraints are enforced (stored in Context)
- Not currently implemented (file-based + in-memory cache only)
- Multi-user support would require middleware layer in `pkg/server/routes.go`
- sync.RWMutex in Server protects concurrent access to trips, dirtyTrips, nameToFile maps
- Background persistence runs in separate goroutine, synchronized via channels (stopPersist, persistDone)
- Atomic creation operations hold lock for entire create sequence (duplicate check + file I/O)
<!-- GSD:architecture-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd:quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd:debug` for investigation and bug fixing
- `/gsd:execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd:profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
