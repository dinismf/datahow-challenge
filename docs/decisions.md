# Notable Decisions

## Architecture

Dependencies flow inward only: `Presentation → Service → Domain`, `Infrastructure → Domain`.

**Domain** holds business entities, repository interfaces (e.g. `FeatureFlag`, `UserOverride`), and generic types like
errors.  
**Service** contains business logic and rules (evaluation logic).  
**Presentation** exposes HTTP handlers, routing, and request/response mapping.  
**Infrastructure** implements domain repository interfaces (in-memory store, future RDBMS/Redis adapters).

**Rationale:** I followed clean architecture principles while deliberately avoiding overengineering for this
demonstration project. The goal was to showcase clean architecture's key benefits — swappable implementations and
isolated testing — without introducing unnecessary abstractions. Each layer can be swapped or tested in isolation;
replacing the in-memory store with Postgres requires only a new infrastructure package. This strikes a balance between
architectural clarity and pragmatism suitable for a project of this scope.

---

## Data modelling

`FeatureFlag` holds only its own state (`id, name, global_enabled, created_at, updated_at`). User overrides were initially embedded as `map[string]bool` on the flag; this was dropped because it doesn't scale (loading a flag pulls all overrides), loses metadata (no timestamps, no audit trail), and doesn't map to any real storage backend.

`UserOverride` is a first-class entity (`flag_id, user_id, enabled, created_at, updated_at`) with its own repository. Would map to a `user_overrides` table in RDBMS or one key per override in Redis (`flag:{id}:override:{user_id}`). `Set` is an upsert, matching `PUT` semantics.

---

## Evaluation order

Asymmetric model — `global=on` always wins; user overrides only apply when global is off:

| global | override   | result   | reason          |
|--------|------------|----------|-----------------|
| on     | any / none | enabled  | `global`        |
| off    | on         | enabled  | `user_override` |
| off    | off / none | disabled | `global`        |

Rationale: The requirements were vague initially, and I first tried implementing user overrides to prevail over global
state, but ended up with this approach instead. Overrides exist to grant early access before a rollout, reduce risk in
deploying new features, test behavior with controlled users in production, and gather feedback before full release.
Once a flag is globally on it is fully released — excluding individual users at that point is not a feature flag
concern.

---

## Error handling

Three layers:

1. **Infrastructure** — storage-specific errors wrap domain errors via `%w` (e.g.
   `fmt.Errorf("memory.FeatureFlagInMemoryRepository: key %q: %w", key, core.ErrNotFound)`). The service uses
   `errors.Is` and never sees the storage type.

2. **Service** — returns `*ServiceError{Code, Message, Reason}`. `Error()` returns only `Message` (client-safe).
   `LogError()` returns the full chain including `Reason` (internal only). `WithReason()` copies the catalog error so
   shared values are never mutated.

3. **Presentation** — a single `httpError` function switches on `Code`, logs `LogError()`, and sends only `Message` to
   the client. Internal errors get a generic message regardless of `Reason`.

---

## Response types in domain

The service returns `domain.FeatureFlagResponse`, `domain.EvaluationResponse`, etc. rather than raw domain entities.
Additionally, some service methods accept raw HTTP request structs directly instead of DTOs, avoiding an extra mapping
layer at the presentation boundary. Both choices trade strict layer purity for removing redundant transformations —
acceptable with a single presentation layer and stable shapes. The tradeoff: slightly tighter coupling between
presentation and service (changes to request shape ripple into service signatures), but simpler code and fewer
allocations. Reconsider if multiple transports need divergent request/response shapes or if the API contract diverges
significantly from domain operations.

---

## Testing strategy

Tests mock at the repository boundary, not the service boundary: `mock repo → real service → real handler → httptest`.
This covers business logic, HTTP contract (status codes, response bodies), and routing in one pass without defining a
service interface solely for testing.

We also have isolated unit tests for the in-memory repository implementation. In an ideal scenario, each layer would be
tested in isolation (unit tests for repositories, services, and handlers independently), and then integration tests
would verify all components working together in a real scenario. The current approach balances thoroughness with
pragmatism for a demonstration project.