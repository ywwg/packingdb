# Agent Instructions

This document provides comprehensive guidance for AI agents working on the packingdb project. For user-facing documentation, see [README.md](README.md) and [WEB_IMPLEMENTATION.md](WEB_IMPLEMENTATION.md).

**⚠️ IMPORTANT: When making changes to the project, always update this document to reflect:**
- New features or functionality added
- Architecture or design pattern changes
- New dependencies or tools introduced
- Updated workflows or best practices
- Common issues discovered and their solutions

## Project Overview

PackingDB is a trip packing list management system with both a TUI (Terminal UI) and web interface. It helps users:
- Create and manage trips with specific contexts (nights, temperature, properties)
- Configure trip properties (Camping, Swimming, Hiking, etc.)
- Track what items are packed for each trip
- Generate packing lists based on trip context

### Architecture

**Core Library (`pkg/packinglib/`):**
- `registry.go` - Property and item registry
- `trip.go` - Trip data structures and operations
- `context.go` - Trip context (nights, temp, properties)
- `item.go` - Item definitions and packing state

**Web Backend (`cmd/packingweb/`):**
- `main.go` - HTTP server entry point with command-line flags
- Uses `pkg/api/routes.go` for REST API handlers
- Uses chi router for clean route definitions

**Web Frontend (`static/`):**
- `index.html` - Alpine.js-based single-page application with Tailwind CSS (274 lines)
- `app.js` - Alpine.js reactive state management (250 lines)
- `styles.css` - Minimal custom CSS (4 lines, only x-cloak directive)

**TUI (`cmd/packingdb/`):**
- Terminal-based interactive interface using promptui

### Key Design Decisions

1. **File-based storage** - Trips stored as YAML/CSV files in `public/trips/`
2. **No authentication** - Single-user system, no login required
3. **Mobile-first** - Web UI designed primarily for phones
4. **Alpine.js framework** - Lightweight reactive framework, no build step
5. **Tailwind CSS** - Utility-first CSS framework via CDN, minimal custom CSS
6. **Stateless server** - Each API request is independent, trips cached in memory
7. **Chi router** - Clean, standard Go routing library
8. **CDN-based frontend** - No build tools, npm, or node dependencies required

## Server Management

### After Manual Testing
Always kill the test server after completing manual testing to avoid:
- Port conflicts on subsequent runs
- Resource leaks
- Confusion about which server is running

**Command:**
```bash
pkill -f "go run ./cmd/packingweb"
```

Or if the binary was built:
```bash
pkill packingweb
```

### Starting the Server for Testing
From repo root:
```bash
go run ./cmd/packingweb
```

With custom flags:
```bash
go run ./cmd/packingweb -port 9000 -trips ./custom/trips -static ./custom/static
```

### Verifying Server is Running
```bash
pgrep -la packingweb
```

Or check the port:
```bash
lsof -ti:8080
```

## Testing Guidelines

### E2E Tests
- E2E tests automatically start/stop their own server with random ports
- No manual cleanup needed after running `go test`
- Tests use temporary directories that are automatically cleaned up
- Located in `cmd/packingweb/e2e_test.go`

**Run E2E tests:**
```bash
cd cmd/packingweb
go test -v -run TestE2E
```

### Manual Testing
- Start server with `go run ./cmd/packingweb`
- Test functionality in browser or with curl
- **ALWAYS** kill the server when done with `pkill -f "go run ./cmd/packingweb"`

### API Testing with curl
```bash
# List trips
curl -s http://localhost:8080/api/trips | python3 -m json.tool

# Get trip details
curl -s http://localhost:8080/api/trips/test-trip | python3 -m json.tool

# Get properties
curl -s http://localhost:8080/api/properties | python3 -m json.tool

# Toggle item
curl -s -X POST http://localhost:8080/api/trips/test-trip/items/z/toggle
```

## Development Workflow

1. Make changes to code
2. Start server for testing: `go run ./cmd/packingweb`
3. Test changes (browser or curl)
4. **Kill server:** `pkill -f "go run ./cmd/packingweb"`
5. Run automated tests if needed: `go test ./...`
6. Commit changes

## REST API Reference

See [WEB_IMPLEMENTATION.md](WEB_IMPLEMENTATION.md#rest-api-endpoints) for complete API documentation.

**Key endpoints:**
- `GET /api/trips` - List all trips
- `POST /api/trips` - Create new trip
- `GET /api/trips/{name}` - Get trip details
- `PUT /api/trips/{name}/update` - Update trip settings (name is stored in file, filename stays the same)
- `GET /api/properties` - List all available properties
- `GET /api/trips/{name}/properties` - Get trip properties
- `POST /api/trips/{name}/properties/{property}/toggle` - Toggle property
- `GET /api/trips/{name}/items` - Get packing items
- `POST /api/trips/{name}/items/{code}/toggle` - Toggle item packed status

**Important behaviors:**
- Toggling a property may automatically enable other properties (e.g., enabling "Camping" might auto-enable "Outdoors")
- After toggling properties, always fetch the full property list to see cascading changes
- Items are organized by category in the response
- Packed status is persisted to the trip file immediately
- Trip name can be changed independently of the filename (stored inside the file)

## Frontend Development

### Alpine.js Patterns

The frontend uses Alpine.js for reactive state management. Key patterns:

**State management:**
```javascript
function packingApp() {
    return {
        currentPage: 'main-menu',
        trips: [],
        currentTrip: null,
        collapsedCategories: new Set(), // Track collapsed state
        // ... other state

        init() {
            this.loadTrips();
        }
    }
}
```

**Template directives:**
- `x-data="packingApp()"` - Initialize Alpine component
- `x-show="condition"` - Conditionally show elements
- `x-for="item in items"` - Loop over arrays
- `x-model="property"` - Two-way data binding
- `@click="method()"` - Event handlers
- `x-text="expression"` - Dynamic text content
- `:class="{ 'active': isActive }"` - Dynamic classes
- `x-collapse` - Smooth collapse/expand animations

**Computed properties:**
Use getters for derived state:
```javascript
get filteredProperties() {
    return this.properties.filter(p => p.name.includes(this.search));
}
```

**Collapsible sections:**
```javascript
toggleCategoryCollapse(categoryName) {
    if (this.collapsedCategories.has(categoryName)) {
        this.collapsedCategories.delete(categoryName);
    } else {
        this.collapsedCategories.add(categoryName);
    }
}

isCategoryCollapsed(categoryName) {
    return this.collapsedCategories.has(categoryName);
}
```

**Auto-refresh for multi-device usage:**
```javascript
startAutoRefresh() {
    this.stopAutoRefresh();
    // Refresh every 10 seconds
    this.refreshInterval = setInterval(() => {
        if (this.currentPage === 'packing') {
            this.refreshItems();
        }
    }, 10000);
}

stopAutoRefresh() {
    if (this.refreshInterval) {
        clearInterval(this.refreshInterval);
        this.refreshInterval = null;
    }
}
```
The packing page automatically refreshes every 10 seconds to pick up changes made in other browsers or devices. The refresh is silent (no loading spinner) and preserves UI state like hidePacked and collapsed categories.

### Styling Guidelines

- Mobile-first approach (design for phones, scale up)
- Minimum tap target: 44px (Tailwind: `h-11` or larger)
- Use gradient theme: `#667eea` to `#764ba2`
- Touch-friendly animations (scale on tap)
- Support safe areas for notched devices
- Tailwind utility classes for all styling
- Minimal custom CSS (only framework-required styles like x-cloak)
- Responsive breakpoints: sm (640px), md (768px), lg (1024px)

### State Synchronization

After API mutations, always refresh state to catch server-side changes:
- After toggling property: fetch all properties
- After toggling item: optimistic update is OK (no cascading changes)
- After updating trip: reload trip details
- After creating trip: reload trip list

**Multi-device support:**
- Packing page auto-refreshes every 10 seconds to detect changes from other browsers/devices
- Refresh is silent (no loading indicator) and preserves UI state (hidePacked, collapsed categories)
- Auto-refresh stops when navigating away from packing page

## Common Issues

### Port Already in Use
If you get "bind: address already in use":
```bash
lsof -ti:8080 | xargs kill -9
```

### Finding Running Servers
```bash
pgrep -la packingweb
# or
ps aux | grep packingweb
```

### Static Files Not Loading
- Check that `-static` flag points to correct directory
- Default is `./static` from repo root
- From `cmd/packingweb` directory, server auto-detects wrong path

### Properties Not Updating After Toggle
- Always reload full property list after toggle
- Properties can have cascading dependencies
- Don't rely on optimistic updates for properties

## Code Organization

```
packingdb/
├── cmd/
│   ├── packingdb/          # TUI application
│   └── packingweb/         # Web server
│       ├── main.go         # Entry point, CLI flags
│       └── e2e_test.go     # End-to-end tests
├── pkg/
│   ├── api/
│   │   └── routes.go       # REST API handlers, chi router
│   ├── contexts/
│   │   └── registry.go     # Property definitions
│   └── packinglib/
│       ├── registry.go     # Core registry interface
│       ├── trip.go         # Trip data structures
│       ├── context.go      # Trip context
│       └── item.go         # Item definitions
├── static/
│   ├── index.html          # Alpine.js SPA
│   ├── app.js              # Alpine.js app logic
│   └── styles.css          # Mobile-first CSS
└── public/
    └── trips/              # Trip YAML/CSV files
```

## Dependencies

### Backend
- Go 1.22+ (uses Go modules)
- `github.com/go-chi/chi/v5` - HTTP router
- Standard library for everything else

### Frontend
- Alpine.js 3.x (CDN) - Reactive framework
- Tailwind CSS 3.x (CDN) - Utility-first CSS framework
- No build tools required
- No npm/node dependencies

## File Formats

### Trip Files (YAML)
Trips are stored as YAML with this structure:
```yaml
name: weekend-camping
nights: 2
temperatureMin: 50
temperatureMax: 75
properties:
  - Camping
  - Swimming
items:
  - code: z
    packed: true
  - code: aa
    packed: false
```

### CSV Support
Legacy CSV format also supported for backwards compatibility.

## Performance Considerations

- Trips are cached in memory after first load
- Trip cache invalidated on any modification
- No database - all file I/O is synchronous
- Suitable for single-user, low-volume usage
- For high-volume, consider adding proper caching layer

## Future Development Notes

Potential improvements documented in [WEB_IMPLEMENTATION.md](WEB_IMPLEMENTATION.md#future-enhancements):
- User authentication/multi-user support
- Offline mode with service workers
- Push notifications for packing reminders
- Import/export functionality
- Trip sharing
- Weather API integration

## Related Documentation

- [README.md](README.md) - User-facing quick start guide
- [WEB_IMPLEMENTATION.md](WEB_IMPLEMENTATION.md) - Detailed web implementation docs
- [static/README.md](static/README.md) - Web interface user guide
