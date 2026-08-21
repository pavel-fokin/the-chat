# The Chat

One limitless conversation with LLM models — no "New Chat" button, no thread
list. Context usage is visible, and Clear empties it. Full intent:
`docs/01-product.md`.

## Project status

Early-stage but no longer just a scaffold: a Go backend serves a React frontend as a single embedded binary, and a domain layer has started to take shape.

## Project principles

The project should respect the following four principles on any layer — code, design, UX, and development process alike.

1. **Cognitive state is structurally bounded.** Only a limited amount of structured information can fit into active focus — both the number of information chunks and the number of abstraction levels.
2. **Processing cognitive state is serial.** Conscious processing and understanding of information is strictly serial.
3. **State is constructed from signals.** Understanding is built only from what's currently in focus — so information should be explicit, not assumed; anything left implicit won't be picked up.
4. **Acquired cognitive structures are unstable.** Once constructed, understanding doesn't stay intact on its own — objects change, relationships weaken, parts become hard to reach. Favor structure (types, names, boundaries) as the durable record over prose that can drift out of sync.

The project should be built to reduce cognitive load at all levels. Less code, words straight to the point, less documentation, but not trade for correctness.

## Doc router

Nothing under docs/ loads on its own. Open a file when its trigger fires:

| Trigger | Read |
|---|---|
| understanding or revisiting the process itself | `docs/00-process.md` |
| the issue is about product or job to be done | `docs/01-product.md` |
| the task touches domain model, architecture or decision | `docs/02-system.md` |
| you need to know what the product does now | `docs/03-features.md` |

Rules under `.claude/rules/` are absent from this table on purpose. They carry
`paths:` frontmatter and load themselves when you open a file they match, so
there is nothing to decide.

## Planning vs building

Plan in a separate session, in plan mode (`Shift+Tab` or `/plan`).

- A planning session ends with Linear issues that carry acceptance criteria.
  It does not end with a list of steps.
- Ask which decisions the work forces and what each costs. Don't ask for a
  step-by-step walkthrough.
- A decision that came up while planning becomes an ADR before building starts.
- If the plan runs past one screen, the issue is too big. Split the issue.

Each issue is one vertical slice: the smallest change that runs end to end
and can be watched working. Never split by layer — a model with no caller
can't be verified on its own.

- Can't write the criteria before digging in? The issue is too big.
- A new domain concept isn't an issue. It lands as an ADR first, before any
  slice uses it.
- A test isn't an issue. It belongs in a slice's acceptance criteria.

## Decisions and rules

Three records, three jobs. Keep them apart.

**An ADR is a system-level decision that's hard to revert.** It earns an entry
when reversing it would touch code you haven't written yet. A choice contained in
one file is a comment, not a decision. Decisions taken together about one thing
are one entry. ADRs are numbered `###` entries under `docs/02-system.md`
§Decisions, append-only — the template and conventions live in
`.claude/rules/adr.md` and load when you open that file.

**A product constraint is something the product won't do.** It lives in
`docs/01-product.md` §Constraints, never in an ADR. If the reasoning is about
what people need rather than what the code can bear, it belongs there.

**A rule is the live constraint either one creates.** It sits under
`.claude/rules/`, carries `paths:` frontmatter, and loads when you open a file it
governs. It says what to do, not why.

**A rule and its source link both ways.** The rule names its source — an ADR
number or a product constraint — and an ADR names its rule in the `Rule:` field.
The source holds the reasoning you'd otherwise re-litigate; the rule is what
stops you violating it without noticing. A rule pointing at a deleted or
superseded source is a bug.

Rules are context, not enforcement, and they load on read rather than on write.
See `docs/00-process.md` §What earns an ADR for what that does and doesn't buy.

## Report your decisions

End every session with a list of the calls you made along the way, one line
each. Not a summary of the work — the decisions inside it.

When a choice is yours to make, make it, then say what you chose. Don't ask
about every small thing, and don't stay quiet about it either.

## Session boundaries

- One session, one Linear issue, one branch.
- Spot an unrelated problem? Don't fix it here. File a separate issue and
  stay on the current task.
- Never auto-approve creating or editing an issue. Always confirm.

Working several issues at once is fine, one worktree each:
`claude --worktree feature-name`. Two limits apply. Only one in-flight issue
may touch `docs/`, since worktrees isolate code but not the shared docs.
And claim the ADR number when you open the issue, not when you write the entry.

## Architecture

One binary: the React frontend is compiled into the Go server via `go:embed`.

| Path | Holds |
|---|---|
| `web/` | React 19 + TypeScript, Vite → `web/dist` |
| `web/web.go` | `//go:embed all:dist` → `web.Dist embed.FS` |
| `internal/server` | `New(fs.FS) http.Handler`, `index.html` fallback for client routes |
| `internal/thechat` | domain entities (`User`, `ID[T]`) — no HTTP, no storage |
| `cmd/the-chat-server` | wiring + `http.ListenAndServe` |

**The frontend must be built before `go build`** — `go:embed` reads `web/dist` at
compile time. `make build` orders them. No dev-mode split yet.

Why any of it is shaped this way: `docs/02-system.md` §Decisions.

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

## Writing

Docs are read under load by someone rebuilding context from zero. Write for that
reader.

- One idea per sentence. Average under 20 words, hard cap around 30.
- Say the thing, then explain it. Never the reverse.
- Active voice, present tense. Name the actor.
- No hedging. Cut "probably", "it might be worth", "in some cases". If something
  is genuinely uncertain, say what would settle it.
- Define a term before using it, or link to where it's defined. Principle 3: a
  reader who doesn't already hold the concept won't pick it up from context.
- Cut adjectives that carry no constraint: "robust", "clean", "simple",
  "powerful", "seamless".
- Parallel content goes in a table or list, not a paragraph.
- Delete a sentence before rewriting it. Most don't survive the question "what
  breaks if this is gone?"

Target Flesch Reading Ease 50 or above, Flesch-Kincaid grade 12 or below. Treat
it as a smoke test, not a gate — technical terms drag the score down honestly.
If a doc scores below, look at sentence length and clause nesting first, never at
vocabulary. Simplifying the words instead of the sentences produces vague prose
with a good score, which is the failure this rule exists to prevent.

## Commit messages

Conventional commits, subject line under 72 chars (e.g. `feat: add UserID with UUIDv7`, `fix: restore dark mode on system theme`). Common prefixes: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `build`, `ci`, `perf`.

## Known gaps

No tests (Go or frontend), no CI, no `LICENSE`. `cmd/the-chat-server/main.go` uses raw `http.ListenAndServe` with no timeouts or graceful shutdown — fine for local dev, worth revisiting before deploying anywhere real.
