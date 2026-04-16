# Codebase Structure

**Analysis Date:** 2026-04-16

## Directory Layout

```
packingdb/
├── cmd/                    # Entry point binaries
│   ├── packingdb/          # Interactive TUI (promptui-based) - primary CLI
│   ├── packingcli/         # Basic CLI (older non-interactive)
│   └── packingweb/         # HTTP server for web UI
├── pkg/                    # Shared libraries
│   ├── contexts/           # Pre-built trip contexts and item registry population
│   ├── items/              # Item definitions organized by category
│   ├── logger/             # Structured logging utilities
│   ├── packinglib/         # Core domain logic (Trip, Item, Context, Registry)
│   ├── server/             # REST API implementation and persistence
│   └── weather/            # Weather/temperature lookup service
├── static/                 # Web UI assets (HTML, CSS, JS)
├── .planning/              # Planning and analysis documents
├── go.mod, go.sum          # Go module dependencies
├── README.md               # Project overview
├── TODO.md                 # Implementation backlog
├── WEB_IMPLEMENTATION.md   # Web feature roadmap
└── AGENTS.md               # Agent workflow instructions
```

## Directory Purposes

**cmd/:**
- Purpose: Binary entry points, not importable as libraries
- Contains: main() functions with flag parsing and server/CLI initialization
- Each subdirectory produces one executable binary

**cmd/packingweb/:**
- Purpose: Web server entry point
- Key files: `main.go` (server initialization, HTTP setup), `e2e_test.go` (integration tests)
- Serves: Static files from `static/` directory, REST API via chi router

**cmd/packingdb/:**
- Purpose: Interactive terminal UI for trip management
- Key files: `packingdb.go` (menu system, main loop, trip orchestration)
- Uses: promptui library for multi-choice prompts and text input

**cmd/packingcli/:**
- Purpose: Simple non-interactive CLI (legacy)
- Key files: `packingcli.go`
- Status: Superseded by packingdb but kept for backward compatibility

**pkg/packinglib/:**
- Purpose: Core domain logic independent of any transport layer
- Contains:
  - `trip.go`: Trip orchestration, item filtering, packing state
  - `context.go`: Context (temperature, nights, properties) definition and management
  - `item.go`: Item definition with prerequisite matching logic
  - `property.go`: Property definitions and PropertySet type
  - `category.go`: Item categorization and sorting
  - `registry.go`: Registry interface and in-memory implementation
  - `yaml.go`: YAML serialization (Trip↔YamlTrip conversion)
  - `packmenu.go`: Menu item representation for TUI
- Key abstraction: All operations preserve immutability of Context/Item while allowing Trip mutations

**pkg/contexts/:**
- Purpose: Registry population with predefined contexts and items
- Key files: `contexts.go` (only file, contains PopulateRegistry() and pre-built context templates)
- Pre-built contexts: Firefly (festival), Cape (beach), Tiny House Spring/Summer/Fall
- Pattern: Direct Context struct initialization + calls to registry methods

**pkg/items/:**
- Purpose: Item definitions organized by semantic category
- Files by category:
  - `art.go`: Art supplies and performance gear
  - `business.go`: Work/professional items
  - `clothing.go`: General clothing (socks, underwear, outerwear)
  - `electrical.go`: Chargers, cables, electronics
  - `entertainment.go`: Books, games, media
  - `food.go`: Snacks, meal prep items
  - `performance.go`: DJ/music equipment
  - `sport.go`: Athletic gear
  - `swim.go`: Swimwear and water items
  - `tasks.go`: Task-specific items (cleaning supplies, handy items)
  - `toiletries.go`: Personal hygiene items
  - `tools.go`: Tools and hardware
  - `unmentionables.go`: Undergarments and intimate items
  - `camp{stuff,ping}.go`: Camping-specific gear
  - `dog.go`: Pet supplies
  - `documents.go`: IDs, tickets, papers
  - `entertainment.go`: Devices and entertainment
  - `food.go`: Food items
  - `shop.go`: Retail/shopping related items
  - `items.go`: Item factory function (unused?)
- Pattern: Each category file exports a function that returns []*Item with appropriate prerequisites

**pkg/server/:**
- Purpose: HTTP API implementation and background persistence
- Key files:
  - `server.go`: Server struct, trip loading, file discovery, name→path mapping
  - `routes.go`: HTTP handler functions (list, create, get, update trips; toggle properties/items)
  - `persist.go`: Background persistence goroutine and dirty trip tracking
- Pattern: chi router mounted at `/api/` prefix; JSON request/response bodies
- Concurrency: sync.RWMutex protects shared trip cache and file mappings

**pkg/logger/:**
- Purpose: Structured logging with level support
- Key files: `logger.go` (single file with Init, Debug, Info, Warn, Error, Fatal)
- Configuration: Set via `-log-level` flag at startup

**pkg/weather/:**
- Purpose: Temperature and weather lookups
- Key files: `weather.go` (WeatherLookup interface, implementation)
- Usage: Currently integrated for location-based weather API calls

**static/:**
- Purpose: Web UI assets served by packingweb
- Files:
  - `index.html`: Main page, responsive mobile layout
  - `app.js`: JavaScript client (API calls, DOM manipulation, state management)
  - `styles.css`: Mobile-first CSS styling
  - `api-examples.html`: Reference documentation for API endpoints
  - `README.md`: Feature documentation
- Pattern: Client-side app communicates with Go server via REST API

## Key File Locations

**Entry Points:**
- `cmd/packingweb/main.go`: Web server - initializes registry, creates Server, starts HTTP listener
- `cmd/packingdb/packingdb.go`: TUI CLI - promptui-based menu system
- `cmd/packingcli/packingcli.go`: Basic CLI - simple non-interactive interface

**Configuration:**
- `.planning/`: Analysis and planning documents
- `go.mod`: Module definition and dependencies
- `AGENTS.md`: Agent workflow instructions
- `TODO.md`: Implementation backlog
- `WEB_IMPLEMENTATION.md`: Web feature roadmap

**Core Logic:**
- `pkg/packinglib/trip.go`: Trip orchestration, list generation, filtering
- `pkg/packinglib/context.go`: Context definition and property management
- `pkg/packinglib/item.go`: Item definition and prerequisite matching
- `pkg/packinglib/registry.go`: Plugin/extensibility registry

**Testing:**
- `pkg/packinglib/item_test.go`: Item prerequisite matching tests
- `pkg/packinglib/yaml_test.go`: YAML serialization round-trip tests
- `pkg/weather/weather_test.go`: Weather lookup tests
- `cmd/packingweb/e2e_test.go`: End-to-end API integration tests

**Data Storage:**
- `public/trips/`: Directory for saved trip files (created at runtime, not in repo)
- Trip files: `<trip-name>-<random-suffix>.yml` (primary) and `.csv` (backup)

## Naming Conventions

**Files:**
- Package main files: `main.go` (entry points only)
- Package other files: `<concept>.go` (trip.go, context.go, item.go)
- Category files in items/: `<category>.go` (clothing.go, electrical.go)
- Tests: `<module>_test.go` (item_test.go, yaml_test.go)
- Directories: `pkg/<package>` (standard Go project structure)

**Directories:**
- `/cmd/<binary>` for each executable
- `/pkg/<package>` for importable libraries
- `/static` for web assets
- `/public` for runtime-created files (trips directory)

**Functions & Types:**
- Exported: PascalCase (NewTrip, Context, Trip, Item)
- Unexported: camelCase (makeList, updateList, extractTripName)
- Constructors: New<Type>() pattern (NewTrip, NewContext, NewItem)
- Getters: simple name or Get<Field>() (Name(), Nights)
- Setters: none (types are mostly immutable or use builder methods)

**Variables:**
- Constants: UPPER_SNAKE_CASE (yaml_version)
- Package vars: camelCase (fireflyContext, capeContext)
- Local vars: camelCase
- Receiver: single letter (t for Trip, c for Context, i for Item, r for Registry)

**Types:**
- Structs: PascalCase (Trip, Context, Item, Server)
- Interfaces: PascalCase (Registry)
- Custom primitives: PascalCase (Category, Property)

## Where to Add New Code

**New Feature (e.g., weather-aware item suggestions):**
- Primary code: `pkg/packinglib/trip.go` (add method to Trip struct)
- API exposure: `pkg/server/routes.go` (add handler function, register with chi router)
- Frontend: `static/app.js` (add JavaScript to call new endpoint and update UI)
- Tests: `pkg/packinglib/trip_test.go` (unit tests) + `cmd/packingweb/e2e_test.go` (integration test)

**New Item Category (e.g., "skiing"):**
- Item definitions: Create `pkg/items/skiing.go` with item list
- Registry population: Update `pkg/contexts/contexts.go` in PopulateRegistry() to register skiing items
- Property definitions: Add new properties to property list in `pkg/packinglib/property.go` if needed

**New Pre-built Context (e.g., "Mountain Winter"):**
- Context definition: Add to `pkg/contexts/contexts.go` as var (e.g., mountainWinterContext)
- Registry registration: Call registry.RegisterContext() in PopulateRegistry()
- Properties: May need to add new Property entries to `pkg/packinglib/property.go`

**New Trip Persistence Format (e.g., JSON):**
- Serialization: Add LoadFromJSON() and SaveToJSON() functions to `pkg/packinglib/trip.go`
- Routing: Update LoadFromFile() and SaveToFile() switch statements (trip.go lines 334–341)
- Tests: Add round-trip tests to ensure format survives load/save cycle

**Utilities/Helpers:**
- Shared helpers: `pkg/packinglib/` (if domain-related) or `pkg/logger/` or `pkg/weather/` (if cross-cutting)
- CLI helpers: `cmd/packingdb/` (menu functions, prompt handling)
- Server helpers: `pkg/server/` (error responses, request parsing)
- Frontend helpers: `static/app.js` (API client functions, DOM utilities)

**Tests:**
- Unit tests: Colocated with source file (`item_test.go` next to `item.go`)
- Integration tests: `cmd/packingweb/e2e_test.go` (spin up real server, make HTTP requests)
- Test data/fixtures: Use builder pattern or factory functions in test files

## Special Directories

**public/trips/:**
- Purpose: Runtime-created directory for persisted trip files
- Generated: Yes (created by packingweb on startup)
- Committed: No (in .gitignore)
- Cleanup: Manual (delete old .yml/.csv files from public/trips/)

**.planning/codebase/:**
- Purpose: Architecture and codebase analysis documents (ARCHITECTURE.md, STRUCTURE.md, etc.)
- Generated: No (manually authored during planning)
- Committed: Yes

**.git/:**
- Purpose: Git repository metadata
- Generated: Yes
- Committed: No (standard .gitignore)

---

*Structure analysis: 2026-04-16*
