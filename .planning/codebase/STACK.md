# Technology Stack

**Analysis Date:** 2026-04-16

## Languages

**Primary:**
- Go 1.25.7 - Backend server and CLI applications (`cmd/packingweb`, `cmd/packingcli`, `cmd/packingdb`)

**Secondary:**
- JavaScript (vanilla) - Frontend logic (`static/app.js`)
- HTML5 - Web interface (`static/index.html`)
- CSS3 - Styling via Tailwind CDN (`static/styles.css`)
- YAML - Trip data serialization (`pkg/packinglib/yaml.go`)

## Runtime

**Environment:**
- Go 1.25.7 runtime

**Package Manager:**
- Go modules (go.mod)
- Lockfile: `go.sum` present

## Frameworks

**Core:**
- chi/v5 v5.2.4 - HTTP router and middleware (`github.com/go-chi/chi/v5`)
  - Used for REST API routing in `cmd/packingweb/main.go`
  - Routes defined in `pkg/server/routes.go`

**Frontend UI:**
- Alpine.js 3.15.8 - Lightweight reactive frontend framework
  - CDN-delivered from jsDelivr
  - Used for state management and reactivity in `static/app.js`
- Tailwind CSS (dynamic CDN build) - Utility-first CSS framework
  - Dynamically compiled from CDN
  - Theme configured in `static/index.html` with primary color extension
  - Used for responsive, mobile-first styling

**Testing:**
- testify v1.9.0 - Assertion library and test utilities (`github.com/stretchr/testify`)
- Go built-in `testing` package

**CLI Interactions:**
- promptui v0.9.0 - Interactive terminal UI for CLI prompts (`github.com/manifoldco/promptui`)
  - Used in `cmd/packingdb/packingdb.go` for TUI-based packing interface

**Build/Dev:**
- Go build tools (no dedicated build system)

## Key Dependencies

**Critical:**
- chi/v5 v5.2.4 - HTTP routing, fundamental to REST API server
- promptui v0.9.0 - Interactive CLI experience for traditional interface
- golang.org/x/crypto v0.22.0 - Cryptographic utilities (in transitive chain)

**Infrastructure:**
- gopkg.in/yaml.v3 v3.0.1 - YAML parsing/serialization for trip data storage
- golang.org/x/sys v0.19.0 - System-level utilities (indirect, via promptui)
- golang.org/x/term v0.19.0 - Terminal utilities (indirect, via promptui)

## Configuration

**Environment:**
- Configured via command-line flags in `cmd/packingweb/main.go`:
  - `-trips` - Directory for trip files (default: `./public/trips`)
  - `-static` - Directory for static files (default: `./static`)
  - `-port` - Server port (default: `8080`)
  - `-log-level` - Logging level: debug, info, warn, error (default: `info`)

**Build:**
- `go.mod` and `go.sum` files
- No build configuration files (Makefile, gradle, etc.)
- Shell scripts provided: `packingdb.sh`, `packingcli.sh` for quick execution

## Platform Requirements

**Development:**
- Go 1.25.7 installed
- GOPATH configuration (noted in README.md as required for module resolution)
- Standard Unix build tools (make, etc.) implicit in shell scripts

**Production:**
- Go 1.25.7 runtime (single binary deployment)
- Directories created at runtime: `./public/trips` (trip data), static file serving from `./static`
- HTTP port exposed (configurable, default 8080)
- Local filesystem access for trip file persistence

**Web Frontend:**
- Modern browser with ES6 support
- JavaScript enabled
- No package manager (CDN-delivered dependencies)

---

*Stack analysis: 2026-04-16*
