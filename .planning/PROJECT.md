# PackingDB: Multi-Trip State Isolation

## What This Is

A packing list manager with CLI and web interfaces. Users create trips with properties (camping, cold weather, etc.), and the system generates a filtered packing list. The web server (`packingweb`) serves a REST API and static frontend. This milestone fixes a state bleed bug and adds proper multi-trip support to the backend.

## Core Value

Each trip's packing state must be completely independent. Packing items in Trip A must never affect Trip B.

## Requirements

### Validated

- ✓ Trip creation with properties, nights, temperature range — existing
- ✓ Item filtering based on context/properties — existing
- ✓ Pack/unpack items via API — existing
- ✓ Background persistence of dirty trips to disk — existing
- ✓ YAML and CSV trip file serialization — existing
- ✓ Multiple entry points (CLI, TUI, web) sharing core library — existing

### Active

- [ ] Server holds multiple trips in memory with isolated packing state
- [ ] Creating or loading a new trip has zero state bleed from other trips
- [ ] Switching between trips via API preserves each trip's pack progress
- [ ] Pack state persists to disk and survives server restart
- [ ] No shared mutable references between trips (items, properties, etc.)

### Out of Scope

- UI changes — the frontend already loads trips by name, no changes needed
- Multi-user/auth support — single user, not relevant to this bug
- Context/property data model refactor — separate concern (noted in CONCERNS.md)
- Panic-to-error migration — separate concern (noted in CONCERNS.md)

## Context

The bug manifests when the packingweb server runs long-term. User packs all items for Trip A, then creates Trip B. Trip B's items appear already packed. The root cause is likely shared item references or registry-level mutation. The `Server.trips` map caches trips keyed by name, and items carry mutable `packed` state on the Item struct.

Key files:
- `pkg/server/server.go` — Server struct with trips map, mutex
- `pkg/server/persist.go` — Background dirty trip persistence
- `pkg/packinglib/trip.go` — Trip struct, makeList(), Pack()
- `pkg/packinglib/item.go` — Item struct with mutable packed state
- `pkg/packinglib/registry.go` — Registry holding item definitions

The fix is architectural: ensure items are deep-copied (or packing state is stored separately) so no trip can mutate another trip's state.

## Constraints

- **Backend only**: No frontend changes. API contract stays the same.
- **Existing persistence**: Dirty trip sync already writes to disk every 30s. Leverage it, don't replace it.
- **Brownfield**: This is an existing Go codebase. Follow existing patterns (chi router, packinglib structure, YAML/CSV serialization).

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Fix via multi-trip isolation, not state reset | Resetting state on trip switch is fragile and doesn't address the root cause of shared references | — Pending |
| Backend only, no UI changes | Frontend already addresses trips by name via API | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-04-16 after initialization*
