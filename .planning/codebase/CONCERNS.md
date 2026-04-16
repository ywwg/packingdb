# Codebase Concerns

**Analysis Date:** 2026-04-16

## Tech Debt

### Data Model Mismatch: Properties vs Contexts

**Issue:** The core data model conflates two distinct concepts—properties (boolean trip attributes like "camping") and contexts (named trip templates like "The Cape"). The code treats contexts as properties but saves only properties to files.

**Files:**
- `pkg/packinglib/context.go` lines 62-77 (commented-out recursive property loading)
- `pkg/packinglib/yaml.go` lines 19, 32-34 (only Properties serialized, not contexts)
- `pkg/packinglib/trip.go` lines 25-29 (trip comment acknowledges this issue)

**Impact:**
- Old packing lists become read-only when the registry changes, since reloading loses context information
- Context names stored as properties pollute the property namespace
- Cannot distinguish between "context constraints" and "property constraints" at load time
- Adds confusing implicit behavior (context registers itself as a property)

**Fix approach:**
- Serialize both contexts and properties in YAML format separately
- Create explicit `Context` struct fields in `YamlTrip` (not just properties)
- Update `LoadFromYAML` to restore context relationships, making old lists editable
- Bump YAML version to v2 and migrate legacy files

---

### Panic-Based Error Handling

**Issue:** The codebase uses `panic()` extensively for runtime errors that should be handled gracefully. This makes the app crash on invalid input or missing data.

**Files:**
- `pkg/packinglib/trip.go` lines 186, 240, 379-390 (pack/load operations panic on bad input)
- `pkg/packinglib/registry.go` (duplicate context/item registration panics)
- `pkg/packinglib/property.go` (property registration contradictions panic)
- `cmd/packingcli/packingcli.go` (multiple panic calls for missing args)

**Current instances:** 17+ panic sites

**Impact:**
- CLI crashes on typos or missing data (e.g., pack nonexistent category name)
- Web server crashes on malformed API requests
- Poor user experience: no chance to recover or fix mistakes
- Difficult to debug without stack traces

**Fix approach:**
- Replace panics with proper error returns throughout `packinglib`
- Update callers to handle errors explicitly
- Raise error in web API layer to return HTTP 400/500 responses
- Add integration tests that verify error handling paths

---

### Unimplemented Custom Function System

**Issue:** The YAML schema supports `custom_func_name` on items but there's no mechanism to register or invoke custom functions.

**Files:**
- `pkg/packinglib/yaml.go` lines 102-105 (stub that returns error)
- `pkg/packinglib/item.go` lines 253-272 (customConsumableMutator type exists)

**Impact:**
- Cannot define reusable custom packing logic
- Dead code path: any YAML with custom function fails to load
- Confuses users who see this schema field with no implementation

**Fix approach:**
- Decide: either remove the schema field or implement a function registry
- If implementing: create a `CustomFunc` registry type, populate in contexts init, use in YAML loader
- If removing: deprecate field, document migration path

---

## Known Bugs

### Name-Based Item Lookup is Fragile

**Issue:** `Trip.Pack()` does case-insensitive name matching as a fallback when code-based lookup fails. This is error-prone when items have similar names.

**Files:**
- `pkg/packinglib/trip.go` lines 150-170 (Pack function)

**Symptoms:**
- If two items differ only in capitalization, both pack when one is intended
- No error if the intended item name is misspelled—just silently doesn't pack
- Message printed to stdout instead of returned as error: "tried to pack nonexistant item, ignoring: %s"

**Trigger:**
- Load a CSV with item name "Sweatshirt" when registry has "sweatshirt" (different case)
- Or typo a name in a packing list file

**Workaround:**
- Use short codes (a, b, c, ...) instead of names when packing
- But CSV format requires names, so file-based packing always hits this fallback

**Fix approach:**
- Remove silent string matching in fallback path
- Require exact name match or return error
- Migrate CSV format to include item codes for robustness

---

### Weather API Timeout Not Gracefully Handled in Web

**Issue:** The weather lookup functions have a 10-second timeout, but there's no fallback if the API is slow or unreachable. The web server has no retry logic.

**Files:**
- `pkg/weather/weather.go` lines 17, 417 (10-second timeout set)
- `cmd/packingweb/main.go` (no weather integration visible in main)

**Impact:**
- If Open-Meteo API is down or slow, weather lookup hangs for 10 seconds then fails
- No user-friendly message; just an API error
- No caching of previous lookups

**Fix approach:**
- Add exponential backoff retry logic with 2-3 attempts
- Cache successful weather lookups in memory for 24 hours
- Return sensible defaults (70°F min/max) if lookup fails
- Add user notification: "Using typical weather for this date/location"

---

## Security Considerations

### CORS Allows All Origins

**Issue:** The web server configures CORS with wildcard origin (`*`).

**Files:**
- `WEB_IMPLEMENTATION.md` line 110 (documented design decision)

**Risk:**
- Any website can make requests to the packingdb server and access all trips
- If packingdb is exposed to a network (not just localhost), anyone can read/modify packing lists
- Sensitive trip data (dates, locations, context) is exposed

**Current mitigation:**
- Documented as "suitable for single-user local use only"
- No authentication system exists

**Recommendations:**
- Keep wildcard CORS only if server listens on localhost:8080
- Add transport-level protection: bind to 127.0.0.1 explicitly (not 0.0.0.0)
- If multiuser support is added (from TODO.md), implement auth and restrict CORS origins
- Document security implications clearly in README

---

### No Input Validation on File Operations

**Issue:** File paths are not validated, allowing potential path traversal attacks.

**Files:**
- `pkg/packinglib/trip.go` lines 437-451 (SaveToFile takes arbitrary filename)
- `cmd/packingweb/main.go` lines 19, 32 (trips directory from flag, created without validation)

**Risk:**
- Web API could be tricked into saving files outside the trips directory
- CLI could overwrite arbitrary system files with packing list data

**Mitigation needed:**
- Validate trip names match safe pattern `[a-zA-Z0-9_-]+`
- Use `filepath.Join` with basename extraction to prevent `../` traversal
- Restrict trips directory to a sandboxed location

---

## Performance Bottlenecks

### Inefficient Property Matching on Every Item Load

**Issue:** `Item.Satisfies()` iterates through all prerequisites for every item every time a context changes. With many items and properties, this is O(items × properties × contexts).

**Files:**
- `pkg/packinglib/item.go` lines 52-83 (Satisfies function)
- `pkg/packinglib/trip.go` lines 121-134 (makeList calls AdjustCount for every item)

**Cause:**
- No caching of item-property satisfaction results
- Mutator sorting happens on every AdjustCount call
- Trip context changes (adding/removing property) rebuilds entire packing list

**Current capacity:**
- Fine for <100 items
- Noticeable slowdown with >500 items
- Web interface would lag if properties changed frequently

**Improvement path:**
- Cache `Satisfies()` results keyed by (item, property_set)
- Sort mutators once at item registration, not per-use
- Add lazy evaluation: only rebuild affected categories when properties change
- Profile with realistic item counts (probably not critical until >1000 items)

---

### Mutator Priority Sorting Repeated Per Item

**Issue:** Every item sorts its mutators on every `AdjustCount()` call, even though mutator order never changes.

**Files:**
- `pkg/packinglib/item.go` lines 94-96 (sort.Slice called every time)

**Impact:**
- Wasteful CPU for large mutator arrays (rare but possible)
- 2-3 extra allocations per item per calculation

**Fix approach:**
- Sort mutators once at item creation time
- Store as pre-sorted slice

---

## Fragile Areas

### Context Initialization and Registry State

**Files:**
- `pkg/packinglib/context.go` lines 27-50 (NewContext modifies global registry)
- `pkg/packinglib/registry.go` (RegisterContext, RegisterProperty have side effects)

**Why fragile:**
- Creating a context registers its name as a property globally
- No transaction support: if property registration fails partway through, registry is inconsistent
- Tests could pollute each other's state if not careful
- Hard to understand what state exists after failed operations

**Safe modification:**
- Create contexts in sequence during init, not in response to user input
- Validate all properties exist before creating a context
- Add rollback mechanism or use immutable context snapshots
- Test that failed context creation leaves registry unchanged

**Test coverage:**
- No tests for failed property registration scenarios
- No tests for duplicate context/property names
- No tests for out-of-order property/context creation

---

### Item Prerequisite Resolution

**Files:**
- `pkg/packinglib/item.go` lines 52-83 (Satisfies logic)
- `pkg/packinglib/context.go` lines 62-77 (commented-out recursive resolution)

**Why fragile:**
- Allow/disallow logic assumes properties are strings (line 67 comment acknowledges issue)
- Disallow prerequisites prevent items completely if ANY disallow is active
- No way to express "require A AND B" or "require A OR B" constraints
- Mixing context names and property names in prerequisites is implicit and error-prone

**Safe modification:**
- Write explicit tests for allow/disallow interactions (especially mixed cases)
- Add helper: verify all prerequisites exist in registry before attaching to item
- Consider explicit constraint language: `require: ["property1", "property2"]` with AND semantics

**Test coverage:**
- `pkg/packinglib/item_test.go` has some tests but missing edge cases
- No tests for disallow preventing items correctly
- No tests for context-as-property behavior

---

### YAML File Format Versioning

**Files:**
- `pkg/packinglib/yaml.go` lines 9-11 (hardcoded yaml_version = 1)
- `pkg/packinglib/trip.go` lines 431-432 (strict version check, rejects mismatches)

**Why fragile:**
- Any schema change breaks old files (no migration logic)
- CSV format has ad-hoc versioning (V2 prefix) separate from YAML
- Two different file formats competing (CSV vs YAML)

**Safe modification:**
- Implement migration layer: detect version, transform to current format
- Document format evolution in SCHEMA.md
- Add test fixtures for each version
- Keep one primary format (recommend YAML)

---

## Scaling Limits

### File-Based Storage

**Issue:** All data stored as YAML/CSV files in `public/trips/` or working directory.

**Files:**
- `pkg/packinglib/trip.go` lines 437-484 (SaveToFile operations)
- `cmd/packingweb/main.go` lines 19, 32, 37 (trips directory)

**Current capacity:**
- Fine for ~100 trips per user
- File I/O latency becomes noticeable with 1000+ trips
- No indexing: searching all trips is O(n)

**Limit:**
- When: >500 concurrent trips being edited
- Breaks: Web API /api/trips list, reload performance

**Scaling path (if needed):**
- Migrate to SQLite for single-machine deployment (keeps file-based feel)
- Or PostgreSQL for multi-machine/cloud
- Add trip metadata index (name, date, context) for fast lookup
- Implement pagination for /api/trips endpoint

---

### Memory Trip Cache

**Issue:** Web server loads all trips into memory at startup.

**Files:**
- `cmd/packingweb/main.go` (trips directory scanned)
- Implied from architecture: no disk I/O shown per-request

**Current capacity:**
- Fine for ~50 active trips in memory
- Each trip is ~5-10KB of JSON

**Scaling path:**
- Lazy load trips on first access
- Implement LRU cache with disk fallback
- Or use database transaction for each request

---

## Dependencies at Risk

### promptui (CLI Dependency)

**Issue:** `github.com/manifoldco/promptui` is only used in CLI tool, not actively maintained.

**Files:**
- `cmd/packingdb/packingdb.go` (all UI uses promptui)
- `go.mod` (dependency listed)

**Risk:**
- No major security issues known, but library hasn't seen updates in years
- If an issue arises, maintainers may not respond
- Keyboard handling (PageUp/PageDown) has known limitations (line 156 TODO)

**Impact:**
- CLI tool could become unmaintainable
- Platform-specific terminal handling may break with OS updates

**Migration plan (if needed):**
- Replace with `github.com/charmbracelet/bubbletea` (actively maintained, modern)
- Or use built-in `bufio.Scanner` + simpler UI (less featureful but functional)

---

### Open-Meteo API Dependency

**Issue:** Weather lookup hardcodes three Open-Meteo API endpoints with no fallback provider.

**Files:**
- `pkg/weather/weather.go` lines 20-24 (baseURL constants)

**Risk:**
- If Open-Meteo changes API, application breaks with no fallback
- No API key = rate limiting could throttle requests if traffic grows
- Service could go offline unpredictably

**Mitigation:**
- Document the dependency clearly
- Add fallback provider (e.g., weatherapi.com, dark sky API)
- Cache results aggressively (24hr minimum)
- Add circuit breaker: if API is down, use defaults

---

## Missing Critical Features

### Multi-User Support

**Issue:** The TODO.md explicitly calls this "first priority" but it's not implemented.

**Files:**
- `TODO.md` lines 3-11 (comprehensive design discussion)
- `pkg/packinglib/trip.go` (no user_id field)
- `cmd/packingweb/main.go` (no auth)

**Problem:**
- Current system assumes single user
- All trips stored in shared directory with no isolation
- No way to know who created or modified a trip
- Web interface is completely open if exposed to network

**Blocks:**
- Cannot share trips with family/team
- Cannot use this in a multi-person household
- Unsafe to expose web UI beyond localhost

**Design notes from TODO:**
- "probably as simple as attaching user IDs to item lists"
- Concern: "this bullshit with name-based items is probably going to bite us"
- Recommendation: keep YAML as initialization, move to SQL eventually

---

### Property/Context Separation

**Issue:** TODO.md line 17-21 identifies this as unresolved.

**Files:**
- `TODO.md` lines 17-21 (explicit design note)
- `pkg/packinglib/property.go` (Property type, no Context type)
- `pkg/packinglib/context.go` (Context as API wrapper, not stored)

**What's missing:**
- Properties should be reusable attributes (e.g., "camping", "cold")
- Contexts should be named bundles (e.g., "winter mountain trip" = {camping, cold, hiking})
- Current system stores only properties, conflates contexts into property namespace

**Impact:**
- Cannot define context templates with default property sets
- Hard to reuse property groupings
- File schema doesn't support context relationships

---

## Test Coverage Gaps

### Item Prerequisite Logic

**What's not tested:**
- Interaction of Allow and Disallow prerequisites when both present
- Disallow prerequisites correctly preventing items
- Context-as-property behavior (item requires context name)
- Edge case: item with no prerequisites

**Files:**
- `pkg/packinglib/item_test.go` (exists but incomplete)

**Risk:**
- Changes to `Item.Satisfies()` could silently break item visibility
- Edge cases in prerequisite logic could go unnoticed
- Hard to refactor with confidence

**Priority:** Medium

---

### File Loading and Saving

**What's not tested:**
- LoadFromCSV with malformed input (missing fields, invalid values)
- SaveToFile then LoadFromFile round-trip consistency
- Version compatibility: loading V2 files into upgraded app
- Edge cases: empty categories, items with count=0

**Files:**
- `pkg/packinglib/trip.go` lines 333-485 (file I/O heavily uses panic)
- `pkg/packinglib/yaml_test.go` (some tests, but limited)

**Risk:**
- Corrupted files could crash the app instead of gracefully degrading
- Old files might lose data on load/save cycle
- Undetected format evolution breaks compatibility

**Priority:** High

---

### Error Handling Paths

**What's not tested:**
- Registry operations with missing properties/contexts
- Invalid user input to CLI prompts
- Network failures in weather API lookup
- Web API with invalid trip names, bad JSON bodies

**Files:**
- No dedicated error handling tests

**Risk:**
- Panics instead of graceful error messages
- Web API returns 500 instead of 400 for user errors
- CLI crashes instead of showing help

**Priority:** High

---

### Web API Integration

**What's not tested:**
- Full request/response cycle for each API endpoint
- Concurrent requests modifying same trip
- Large payloads (100+ items in pack list)
- Edge cases: trip names with special characters, unicode

**Files:**
- `cmd/packingweb/e2e_test.go` (exists but minimal)

**Risk:**
- Race conditions in file writes
- Malformed JSON responses
- API contract drift from implementation

**Priority:** Medium

---

*Concerns audit: 2026-04-16*
