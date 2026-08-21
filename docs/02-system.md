# Domain model, architecture, decisions

## Domain model

`package thechat` (`internal/thechat/`) holds the domain. It has no dependencies
on HTTP, storage, or the frontend.

**User** (`user.go`) — a person using the chat. Currently just an identity.

**ID[T]** (`id.go`) — every entity gets its own ID type from one generic wrapper,
`ID[T any] struct{ uuid.UUID }`, aliased per entity (`type UserID = ID[user]`).
`T` is a phantom type: it exists only to make the compiler reject a `UserID`
passed where another entity's ID belongs. Adding the next entity's ID is one
line. See decision 0001.

**No container entity.** There is no `Chat`, `Conversation`, `Thread`, or
`Session`, and the absence is deliberate. Messages reference `UserID` directly.
The reason is a product constraint, not a modelling one — see `01-product.md`
§Constraints before adding one.

## Architecture

One binary. The frontend is compiled into it.

| Piece | Holds |
|---|---|
| `web/` | React 19 + TypeScript, built by Vite into `web/dist` |
| `web/web.go` (`package web`) | `//go:embed all:dist` → exposes `web.Dist embed.FS` |
| `internal/server` | `New(fs.FS) http.Handler` — serves the files, falls back to `index.html` so client-side routes survive a reload |
| `cmd/the-chat-server` | wires `web.Dist` → `fs.Sub(…, "dist")` → `server.New` → `http.ListenAndServe`, default `0.0.0.0:8080` (`-addr`) |
| `internal/thechat` | the domain, above — not yet reachable from the server |

`server.New` takes an `fs.FS` rather than the `embed.FS` directly, so it can be
pointed at `os.DirFS("web/dist")` if a disk-backed dev mode is ever wanted.

**Build order is a hard constraint:** `go:embed` reads `web/dist` at compile
time, so the frontend must be built before `go build`. `make build` does both in
order. There is no dev-mode split where Go proxies to the Vite dev server.

## Decisions

Numbered entries. The template and the conventions live in
`.claude/rules/adr.md`, which loads whenever this file is opened.

A decision earns an entry when reversing it would touch code that isn't written
yet. A decision that applies more than once also gets a rule under
`.claude/rules/`, linked from its `Rule:` field. See `00-process.md` for why the
two are separate.

### 0001 — Entity identity is a UUIDv7 inside a phantom-typed `ID[T]`

- Status: Accepted, 2026-08-21
- Job: none (structural)
- Terms: `User`, `ID[T]`
- Rule: `.claude/rules/thechat-domain.md`

**Context.** The first entity needed an ID, and whatever it got would be copied
by every entity after it. Two questions at once: what the value is, and what the
type is.

**Decision.** IDs are UUIDv7 from `github.com/google/uuid`, wrapped in one
generic `ID[T any] struct{ uuid.UUID }` and aliased per entity. `T` is a phantom
type — it exists only so the compiler rejects a `UserID` where another entity's
ID belongs. `NewUserID` uses `uuid.Must`, because `NewV7` fails only if the OS
CSPRNG is broken.

**Rejected.** UUIDv4 — random, so it fragments the index and carries no order.
Database integers — force a round trip before an entity exists and leak the
storage engine into the domain. A hand-written ID type per entity — same safety,
N copies of the same code.

**Costs.** UUIDv7 leaks creation time to anyone holding an ID. 16 bytes per key
instead of 8. `uuid.Must` panics rather than returning an error, which is a
deliberate trade of a recoverable path for call-site simplicity.

**Revisit when.** Moving to a store with no native uuid type, where 16-byte keys
land as strings. Or if leaking creation timestamps becomes a privacy problem.

### 0002 — The frontend ships embedded in the Go binary

- Status: Accepted, 2026-08-21
- Job: none (structural)
- Terms: none
- Rule: none — it constrains the build, not the code

**Context.** A React frontend and a Go server have to reach production somehow.
The project runs on one box with one contributor.

**Decision.** `web/dist` is compiled into the binary with `//go:embed all:dist`.
Deployment is one artifact — no static host, no CDN, no nginx. `web/web.go` sits
inside `web/` because `go:embed` patterns cannot reach upward past their own
directory, so the embedding code must sit next to what it embeds.

**Rejected.** Serving `web/dist` from disk — two things to deploy and keep in
sync. A separate static host — an extra service and a CORS boundary for a
single-box project.

**Costs.** No hot reload: every frontend change needs a full `make build`, and
`go build` alone silently uses a stale `dist`. Binary size grows with the
frontend. Frontend and backend can no longer ship independently.

**Revisit when.** The rebuild loop gets slow enough to interrupt frontend work —
the fix is a dev-mode split where Go proxies to Vite, which `server.New`'s
`fs.FS` parameter already allows. Or when the frontend needs its own deploy
cadence.
