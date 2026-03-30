# datahow-challenge

## Running

```bash
make run    # start the server on :8080
make build  # compile to bin/api
make test   # run all tests
```

## Endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/flags` | Create a flag |
| `GET` | `/flags/:id` | Get flag by key |
| `PUT` | `/flags/:id/global` | Set global on/off |
| `PUT` | `/flags/:id/users/:user_id` | Set per-user override |
| `GET` | `/flags/:id/users/:user_id/evaluation` | Evaluate flag for a user |

## curl Examples

**Create a flag (globally disabled)**
```bash
curl -s -X POST http://localhost:8080/flags \
  -H "Content-Type: application/json" \
  -d '{"key":"new-dashboard","name":"New Dashboard","global_enabled":false}' | jq
```

**Evaluate before any override — returns global state**
```bash
curl -s http://localhost:8080/flags/new-dashboard/users/user-1/evaluation | jq
# {"enabled":false,"reason":"global"}
```

**Set a per-user override (opt user-1 in while global is still off)**
```bash
curl -s -X PUT http://localhost:8080/flags/new-dashboard/users/user-1 \
  -H "Content-Type: application/json" \
  -d '{"enabled":true}' | jq
# {"flag_id":"new-dashboard","user_id":"user-1","enabled":true,"created_at":"...","updated_at":"..."}
```

**Evaluate after override — user override wins**
```bash
curl -s http://localhost:8080/flags/new-dashboard/users/user-1/evaluation | jq
# {"enabled":true,"reason":"user_override"}
```

**Flip the flag on globally**
```bash
curl -s -X PUT http://localhost:8080/flags/new-dashboard/global \
  -H "Content-Type: application/json" \
  -d '{"enabled":true}'
# HTTP 204 No Content
```

**A different user with no override gets the global state**
```bash
curl -s http://localhost:8080/flags/new-dashboard/users/user-2/evaluation | jq
# {"enabled":true,"reason":"global"}
```

**Opt user-2 out while global is on**
```bash
curl -s -X PUT http://localhost:8080/flags/new-dashboard/users/user-2 \
  -H "Content-Type: application/json" \
  -d '{"enabled":false}' | jq
# {"flag_id":"new-dashboard","user_id":"user-2","enabled":false,"created_at":"...","updated_at":"..."}

curl -s http://localhost:8080/flags/new-dashboard/users/user-2/evaluation | jq
# {"enabled":false,"reason":"user_override"}
```

**Get a flag that does not exist**
```bash
curl -s http://localhost:8080/flags/unknown | jq
# {"code":"4000","message":"resource not found"}
```

## Design

Feature flagging decouples deploying code from releasing it — functionality can be gated or rolled out at runtime without a deployment.

The project uses a clean layered architecture (`domain` → `service` → `presentation`, with `infrastructure` implementing `domain` interfaces). Dependencies flow inward only; the in-memory store can be replaced with Postgres or Redis without touching the service or presentation layers.

Notable decisions are documented in **[docs/decisions.md](docs/decisions.md)**, covering:

- Data modelling for flags and user overrides (and why the embedded `map[string]bool` was dropped)
- Error handling strategy across layers (sentinel errors, service error catalog, safe logging vs client responses)
- Testing strategy (sociable unit tests: mock repos → real service → real handler)
- Evaluation order interpretation
- Trade-off: request / response types in `domain` 