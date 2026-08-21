# The Chat

One limitless conversation with LLM models — no "New Chat" button, no thread
list. Context usage is visible, and Clear empties it. Full intent:
`docs/01-product.md`.

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

## Three stages, three sessions

Never plan and build in one session. Each stage starts clean and ends with an
artifact the next stage can read without the transcript.

| Stage | Takes | Produces | Ends when |
|---|---|---|---|
| Pre-planning | a feature idea | ordered issue stubs, title + one line each | every slice has a name and a position |
| Planning | one stub | one Linear issue with criteria, plus doc changes | the criteria are written and the docs are merged |
| Building | one issue | one branch, merged | every criterion passes |

Pre-planning is conditional. Skip it when the feature is already one slice —
splitting a single slice produces nothing. Run it when a feature holds several,
so planning never has to hold more than one at a time.

Plan in plan mode (`Shift+Tab` or `/plan`). Ask which decisions the work forces
and what each costs. Don't ask for a step-by-step walkthrough. If the plan runs
past one screen, the issue is too big — split the issue, not the plan.

Each issue is one vertical slice: the smallest change that runs end to end
and can be watched working. Never split by layer — a model with no caller
can't be verified on its own.

- Can't write the criteria before digging in? The issue is too big.
- A new domain concept isn't an issue. It lands as an ADR first, before any
  slice uses it.
- A test isn't an issue. It belongs in a slice's acceptance criteria.

### What a Linear issue must carry

The issue is the handoff. A building session starts cold and sees only this, so
anything left in the planning transcript is lost. Team: **The Chat**.

- **Title** — one slice, one observable outcome.
- **Job** — which job from `docs/01-product.md` this serves. No job, no issue.
- **Acceptance criteria** — observable conditions, checkable by running the app.
  "Sending a message shows the reply in the thread", not "wire up the handler".
  Written before building starts, or the issue is too big.
- **Applies** — the ADRs, product constraints, and rules that govern this slice.
  Name them, so the building session doesn't rediscover them.
- **Out of scope** — what this slice deliberately doesn't do.

Use the branch name Linear generates for the issue. Move it to In Progress when
building starts, In Review at the PR.

### Which docs change in which stage

| Stage | Doc changes |
|---|---|
| Pre-planning | none |
| Planning | an ADR if a decision surfaced; `01-product.md` if the job or a constraint is new |
| Building | `03-features.md` when the slice ships; `02-system.md` if the domain model changed |

Docs the planning stage owes are merged before building starts. A decision left
in a transcript drifts even while the code stays fine.

### When the process itself gets in the way

File it as a Linear issue with the `process` label, under the same rule as any
other problem you spot mid-session. The first few runs will produce several, and
that is the point — the process is being tested too.

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

- One session, one stage, one artifact. Building is also one branch.
- Spot an unrelated problem? Don't fix it here. File a separate issue and
  stay on the current task — fixing it breaks the boundary, filing it protects
  the boundary.
- Never auto-approve creating or editing an issue. Always confirm.

Working several issues at once is fine, one worktree each:
`claude --worktree feature-name`. Two limits apply. Only one in-flight issue
may touch `docs/`, since worktrees isolate code but not the shared docs.
And claim the ADR number when you open the issue, not when you write the entry.

## Architecture

One binary: the React frontend is compiled into the Go server via `go:embed`.
Layout and rationale are in `docs/02-system.md` — the one thing worth carrying
every session is the build order.

**The frontend must be built before `go build`.** `go:embed` reads `web/dist` at
compile time, so `go build` alone uses a stale `dist`. `make build` orders them.
No dev-mode split yet.

## Commands

- `make install` / `make build` / `make run` — see the `Makefile`
- `go build ./...`, `go vet ./...`, `gofmt -l .` — no `*_test.go` files yet
- `npm --prefix web run dev` — Vite with hot reload, frontend only
- `npm --prefix web run lint` — Oxlint. No test runner configured yet.

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
