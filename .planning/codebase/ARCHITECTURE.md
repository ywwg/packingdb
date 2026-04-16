# Architecture

**Analysis Date:** 2026-04-16

## Pattern Overview

**Overall:** Layered architecture with clear separation between domain logic, API presentation, and data persistence.

**Key Characteristics:**
- Domain-driven design with core packing logic isolated from transport layers
- Multiple entry points (CLI, TUI, HTTP API) all sharing the same core library
- Registry pattern for managing extensible item and property definitions
- Background persistence for server-side changes
- File-based storage with CSV and YAML support

## Layers

**Domain/Core Library (`pkg/packinglib`):**
- Purpose: Pure business logic for managing trips, items, properties, and packing lists
- Location: `pkg/packinglib/`
- Contains: Trip orchestration, item filtering, context management, property sets, file loading/saving
- Depends on: Standard library only (golang.org/x/crypto for password hashing in some cases, gopkg.in/yaml.v3 for serialization)
- Used by: All other packages and entry points

**Registry/Composition Layer (`pkg/contexts`):**
- Purpose: Defines and registers available items, properties, and pre-built contexts
- Location: `pkg/contexts/contexts.go`
- Contains: Item definitions from `pkg/items/*`, property list definitions, context templates (Firefly, Cape, Tiny House variants)
- Depends on: `pkg/packinglib`, `pkg/items`
- Used by: Server startup and CLI initialization

**Item Definition Layer (`pkg/items`):**
- Purpose: Categorized item definitions with property prerequisites
- Location: `pkg/items/` (split by category: clothing.go, electrical.go, toiletries.go, etc.)
- Contains: Item lists with conditional logic based on trip properties
- Depends on: `pkg/packinglib` (Item, Category types)
- Used by: Registry during initialization

**API/Transport Layer (`pkg/server`):**
- Purpose: REST API exposure of domain operations with request/response marshaling
- Location: `pkg/server/`
- Contains: HTTP handlers, routing, JSON serialization, file-to-trip mapping
- Depends on: `pkg/packinglib`, `pkg/logger`, standard http library, go-chi router
- Used by: HTTP server entry point (`cmd/packingweb`)

**CLI/TUI Layer (`cmd/packingdb`):**
- Purpose: Interactive terminal-based interface for trip management
- Location: `cmd/packingdb/packingdb.go`, `cmd/packingcli/packingcli.go`
- Contains: Prompt UI menus, console I/O, trip orchestration via promptui
- Depends on: `pkg/packinglib`, `pkg/contexts`, manifoldco/promptui
- Used by: Terminal users

**Web Server (`cmd/packingweb`):**
- Purpose: HTTP server entry point that initializes registry and starts API
- Location: `cmd/packingweb/main.go`
- Contains: Flag parsing, logger initialization, server creation, graceful shutdown handling
- Depends on: `pkg/server`, `pkg/packinglib`, `pkg/contexts`, `pkg/logger`
- Used by: Users running the web interface

**Utilities:**
- Logger: `pkg/logger/logger.go` - Structured logging with level support
- Weather: `pkg/weather/weather.go` - Temperature/weather lookups (currently supports location-based weather API calls)

## Data Flow

**Create Trip Flow:**
1. User provides trip name, nights, temperature range, and properties via API/CLI
2. `NewContext()` in `pkg/packinglib/context.go` creates a Context with property requirements
3. `NewTripFromCustomContext()` creates a Trip and generates ItemList via `makeList()`
4. Item filtering occurs in `trip.makeList()`: each item from registry is checked against context via `item.Satisfies(context)`
5. Filtered items are stored in Trip.packList (map of Category→Items)
6. Trip is either returned (CLI) or saved to file (API) via `trip.SaveToFile()`

**Update Trip Flow:**
1. API handler receives property toggle or item pack/unpack request
2. Trip is loaded from in-memory cache or disk (LRU via `loadTrip()`)
3. Mutation applied: `trip.AddProperty()` or `trip.ToggleItemPacked()`
4. Trip state change marks trip as "dirty" in `dirtyTrips` map
5. Background goroutine in `pkg/server/persist.go` periodically saves dirty trips to disk (30s interval)

**Load Trip Flow:**
1. File path determined via `extractTripName()` which reads first line of CSV or YAML
2. Extension routed to `LoadFromCSV()` or `LoadFromYAML()`
3. YAML unmarshals to `YamlTrip` struct, then calls `AsTrip()` to convert to domain Trip
4. CSV parses header to extract context, nights, temperatures, and properties
5. Item pack state (true/false, item name) parsed from subsequent CSV lines
6. Trip reconstructed with `NewTripFromCustomContext()`, then packed items are replayed via `trip.Pack()`

**Packing Decision:**
1. Item defines prerequisites as PropertySet (map of Property→bool: allow or disallow)
2. Context provides active Properties
3. Item.Satisfies(context) checks:
   - If ANY disallowed property is active: return false
   - If ANY allowed property is active: return true
   - If only disallowances exist: return true (opt-out model)
   - Otherwise: return false

**State Management:**
- In-memory: Trips cached in Server.trips map (keyed by trip name)
- On-disk: Trips saved to `public/trips/` directory as .yml or .csv files
- Dirty tracking: dirtyTrips map holds names of trips needing persistence
- Background sync: Persistence goroutine runs every 30s, writes only dirty trips

## Key Abstractions

**Trip (`pkg/packinglib/trip.go`):**
- Purpose: Represents a complete packing scenario with context, items, and pack state
- Examples: `cmd/packingweb/main.go` line 37, `pkg/server/routes.go` line 81
- Pattern: Mutable value object that holds filtered ItemList and applies mutations (Pack, AddProperty)

**Context (`pkg/packinglib/context.go`):**
- Purpose: Describes conditions that determine which items are relevant (nights, temperature, properties)
- Examples: Pre-built contexts in `pkg/contexts/contexts.go` (fireflyContext, capeContext, tinyhouseSpring, etc.)
- Pattern: Immutable by design, holds PropertySet that drives item filtering

**Registry (`pkg/packinglib/registry.go`):**
- Purpose: Service locator for all items, properties, and context definitions
- Examples: `pkg/packinglib/registry.go` interface definition, `pkg/contexts/contexts.go` PopulateRegistry()
- Pattern: Interface (Registry) with in-memory implementation (StructRegistry), populated at startup

**Item (`pkg/packinglib/item.go`):**
- Purpose: Individual thing to pack with conditional inclusion based on context
- Examples: `pkg/items/clothing.go`, `pkg/items/toiletries.go` item definitions
- Pattern: Immutable item name/description with mutable packed state and prerequisite matching

**PropertySet (map[Property]bool):**
- Purpose: Efficient set-like structure for property matching in Item.Satisfies()
- Pattern: True=allow, False=disallow; enables flexible opt-in and opt-out filtering

## Entry Points

**Web Server:**
- Location: `cmd/packingweb/main.go`
- Triggers: User runs binary; listens on :8080
- Responsibilities: HTTP server lifecycle, static file serving, API routing, background persistence

**CLI (Interactive):**
- Location: `cmd/packingdb/packingdb.go`
- Triggers: User runs `./packingdb <tripfile.yml>`
- Responsibilities: Interactive menus via promptui, trip manipulation, file I/O

**CLI (Basic):**
- Location: `cmd/packingcli/packingcli.go`
- Triggers: User runs `packingcli` (older interface)
- Responsibilities: Non-interactive command-line trip creation

## Error Handling

**Strategy:** Panic on unrecoverable startup errors, graceful error returns for runtime operations.

**Patterns:**
- File operations return errors (LoadFromFile, SaveToFile)
- Trip creation returns errors (NewTrip, NewContext)
- Property/Item lookup returns errors (AddProperty, Pack)
- HTTP handlers use respondError() with appropriate status codes (400 BadRequest, 409 Conflict, 500 InternalServerError)
- Registry lookups return errors (GetContext, HasProperty)
- Logger.Fatal() used for fatal startup conditions (server creation, directory creation)

## Cross-Cutting Concerns

**Logging:** 
- Framework: `pkg/logger/logger.go` provides Debug, Info, Warn, Error, Fatal functions
- Configured at startup via `-log-level` flag (default: info)
- Used for: Server startup, trip persistence, scan results, API requests

**Validation:**
- Request body validation in HTTP handlers (name required, nights > 0)
- Trip file parsing validates YAML/CSV format and version
- Property existence checked against registry (HasProperty)
- Temperature range constraints are enforced (stored in Context)

**Authentication:**
- Not currently implemented (file-based + in-memory cache only)
- Multi-user support would require middleware layer in `pkg/server/routes.go`

**Concurrency:**
- sync.RWMutex in Server protects concurrent access to trips, dirtyTrips, nameToFile maps
- Background persistence runs in separate goroutine, synchronized via channels (stopPersist, persistDone)
- Atomic creation operations hold lock for entire create sequence (duplicate check + file I/O)

---

*Architecture analysis: 2026-04-16*
