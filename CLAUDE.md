# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project status

The Chat is a project aiming for effective chatting with AI models. It's early-stage but no longer just a scaffold: a Go backend serves a React frontend as a single embedded binary, and a domain layer has started to take shape. Git history starts at an initial commit (`feat: initial structure with go server and react frontend`).

## Project principles

The project should respect the following four principles on any layer — code, design, UX, and development process alike.

1. **Cognitive state is structurally bounded.** Only a limited amount of structured information can fit into active focus — both the number of information chunks and the number of abstraction levels.
2. **Processing cognitive state is serial.** Conscious processing and understanding of information is strictly serial.
3. **State is constructed from signals.** Understanding is built only from what's currently in focus — so information should be explicit, not assumed; anything left implicit won't be picked up.
4. **Acquired cognitive structures are unstable.** Once constructed, understanding doesn't stay intact on its own — objects change, relationships weaken, parts become hard to reach. Favor structure (types, names, boundaries) as the durable record over prose that can drift out of sync.

The project should be built to reduce cognitive load at all levels. Less code, words straight to the point, less documentation, but not trade for correctness.

## Architecture

- `web/` — React 19 + TypeScript frontend, built with Vite. `npm run build` outputs to `web/dist`.
- `web/web.go` (`package web`) — embeds `web/dist` via `//go:embed all:dist` and exposes `web.Dist embed.FS`. This file lives inside `web/`, next to `dist/`, specifically so the embed pattern doesn't need `..`: `go:embed` patterns can only reach downward from the source file's own directory, so the embedding code has to be colocated with (or above) the directory it embeds.
- `internal/server` (`package server`) — `server.New(fs.FS) http.Handler` serves static files from the given filesystem and falls back to `index.html` for any path that doesn't match a real file, so client-side routing (react-router) works on page reloads/deep links. It's decoupled from `embed` — it just takes an `fs.FS`, so it could be pointed at `os.DirFS("web/dist")` for a disk-backed dev mode instead of the embedded copy if that's ever wanted.
- `cmd/the-chat-server/main.go` — entrypoint. Wires `web.Dist` → `fs.Sub(..., "dist")` → `server.New(...)` → `http.ListenAndServe`. Listens on `0.0.0.0:8080` by default (overridable with `-addr`).
- `internal/thechat` (`package thechat`) — domain layer. Holds entities like `User`. Entity IDs use UUIDv7 (`uuid.NewV7()` via `github.com/google/uuid`, the project's first external dependency) for time-ordered, DB-index-friendly IDs. Each entity gets a distinct ID type via a shared generic wrapper, `ID[T any] struct{ uuid.UUID }`, aliased per entity (e.g. `type UserID = ID[user]`) — one definition shared by every entity instead of hand-duplicated types, so a `UserID` can never be passed where another entity's ID is expected, but adding the next entity's ID is one line. `NewUserID()` uses `uuid.Must(uuid.NewV7())`: `NewV7()` only fails if the OS's CSPRNG is broken, an unrecoverable condition, so this fails fast like the rest of the codebase rather than threading an error through every call site.

Because `go:embed` reads `web/dist` at compile time, the frontend must be built *before* `go build`/`go run` — there's no dev-mode split yet where the Go server proxies to the Vite dev server. `make build` handles the ordering.

## Commands

Root (`Makefile`):
- `make build` — `npm --prefix web run build`, then `go build -o bin/the-chat-server ./cmd/the-chat-server`
- `make run` — `make build`, then runs `./bin/the-chat-server`

Go:
- `go build ./...`, `go vet ./...`, `gofmt -l .` — no test files exist yet (`*_test.go`).

Frontend (`web/`):
- `npm run dev` — Vite dev server (bound to `0.0.0.0`)
- `npm run build` — type-check (`tsc -b`) then build with Vite
- `npm run lint` — Oxlint
- `npm run preview` — preview the production build
- No test runner configured yet.

## Web stack notes

- React 19 + TypeScript, bundled with Vite 8 (rolldown-vite) and `@vitejs/plugin-react`. Routing is `react-router-dom` (`BrowserRouter` in `main.tsx`, `Routes`/`Route` in `App.tsx`).
- shadcn/ui is installed (`--template vite -b base`, Base UI primitives, Nova preset): `components.json`, `src/components/ui/`, `src/lib/utils.ts`. The `@` path alias (`./src`) is configured in both `vite.config.ts` and the `tsconfig.*.json` files' `paths` (no `baseUrl` — it's deprecated under this TS version's `moduleResolution: "bundler"`, and `paths` resolves relative to the tsconfig without it).
- Tailwind CSS v4 via `@tailwindcss/vite`, imported in `src/index.css`.
- `src/index.css` carries two token systems that both need to stay in sync: the site's own tokens (`--text`, `--text-h`, `--bg`, `--code-bg`) and shadcn's tokens (`--background`, `--foreground`, `--border`, `--accent`, etc., which shadcn normally toggles via a `.dark` class). Since this project has no theme-toggle mechanism, both sets are switched together by the same plain `@media (prefers-color-scheme: dark)` block — the `.dark` class's values are duplicated into it so shadcn components also follow the OS theme automatically. If a manual toggle is ever added, that duplication should be reconsidered.
- Linting is via Oxlint (`web/.oxlintrc.json`), not ESLint. Type-aware lint rules are not enabled by default; see `web/README.md` for how to add them via `oxlint-tsgolint` if needed.
- `tsconfig.json` is a project-references root pointing at `tsconfig.app.json` (app code) and `tsconfig.node.json` (Vite config).

## Commit messages

Conventional commits, subject line under 72 chars (e.g. `feat: add UserID with UUIDv7`, `fix: restore dark mode on system theme`). Common prefixes: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `build`, `ci`, `perf`.

## Known gaps

No tests (Go or frontend), no CI, no `LICENSE`. `main.go` uses raw `http.ListenAndServe` with no timeouts or graceful shutdown — fine for local dev, worth revisiting before deploying anywhere real.
