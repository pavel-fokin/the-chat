# How this project runs, and why

Read this when you are rethinking the process, not during a working session.
The day-to-day rules live in `CLAUDE.md`, which is short and loads on its own.

## What the process assumes

The whole design rests on one claim: **active state is bounded**, for people
and for models alike. Everything else sits in files. But a file is not
understanding. Understanding gets rebuilt from scratch, from signals, every
time someone starts work.

Four consequences follow, and they shape everything below.

### 1. State is bounded, so context loads in layers

The goal is not to make all knowledge reachable. It is to put the right
knowledge in front of you at the moment you need it. Four layers, sorted by
how they arrive:

| Layer | Holds | Arrives |
|---|---|---|
| `CLAUDE.md` | commands, conventions, the router | every session, automatically |
| `.claude/rules/*.md` with `paths:` | rules for one kind of file | when a matching file opens |
| skills (`/name`) | multi-step procedures | on call, or on a description match |
| `docs/**` | rationale, vocabulary, architecture | when a person decides, via the router |

Layer 2 is in use: `.claude/rules/adr.md` and `.claude/rules/thechat-domain.md`.
Layer 3 is empty — no skills yet. The table describes the scheme, not an
inventory.

The rule of thumb runs bottom-up: the more often you need something, the
higher it belongs. Keep `CLAUDE.md` short. Anything you don't need every
single session drops a layer.

### 2. State changes in sequence, so session boundaries stay hard

One session, one issue, one branch. Every change you accept becomes the
ground for the next choice. Mixing two tasks in one session desynchronizes
the boundary: the context now holds two events, and reasoning about one
bleeds into the other.

The same logic explains why an unrelated bug gets filed rather than fixed on
the spot. Filing it protects the boundary. Fixing it breaks the boundary.

### 3. State is built from signals, so docs must stand alone

A document does not hand over understanding. It offers a signal that the
reader assembles into understanding — and the two readers here assemble
differently. Claude starts every session from zero. You start with a prior
picture that has been quietly drifting.

That sets a hard requirement for `02-system.md`: a definition has to work
without setup. "An episode is a reference over the log" only lands for
someone who already holds the container-versus-view distinction. So the
distinction goes in the entry, instead of being assumed.

### 4. Structures decay, so catching drift needs a machine

Understanding falls apart on its own, even when the documents behind it
never changed. That single fact sets the sharpest constraint on this whole
design:

**Anything that rests on the phrase "update it in the same session" is an
intention, not a mechanism, and it will rot.**

The discipline has to live in artifacts instead: a checklist that fires on an
event, a check that runs in CI, a review tied to a PR that touches structure.
Anything tied to a day of the week is the first habit to go.

**None of that exists yet, and this is the known gap in the process.** Right now
every doc in this directory stays current only because someone remembers, which
is exactly the mechanism this section argues doesn't work. The gap is tolerable
at four documents and one contributor. It stops being tolerable roughly when the
docs get big enough that nobody rereads them end to end.

Three checks are already specified by the rest of these docs, and none of them
runs. They are the concrete first thing to build here:

- A rule whose `Rule:`/ADR backlinks don't resolve, or that cites a superseded
  ADR, is a bug that nothing currently catches.
- Domain terms named in `02-system.md` should appear in the code. Nothing checks
  that they still do.
- The readability target in `CLAUDE.md` §Writing is a number nobody measures.

Until then, each is an intention. Naming them is not the same as enforcing them,
and this paragraph is not a mechanism either.

---

## Planning is a separate session from building

Writing a plan is cheap for the model. Reading one is expensive for you.
Nothing bounds the output, so a forty-step plan arrives in seconds and gets
approved without being read. An unread plan is not a plan. It is a commitment
you never actually made, and afterwards you cannot spot the drift, because
catching drift means holding the plan in active state and you never did.

There is a deeper reason long plans fall apart. Every step you finish changes
the state the next choice gets made from. Step eight was planned against a
state that step three will have destroyed. So the planning horizon is set by
how far the state stays predictable, not by how big the task is, and that is
usually a short way out. Anything past that horizon is fiction, and it drifts
on schedule.

### Plan for issues and criteria, not for steps

Steps describe *how*. They depend on the current state, they rot first, and
the model regenerates better ones once it is actually in the code.

Criteria describe *what has to become true*. They hold regardless of state,
they survive execution, and they can be checked.

So a planning session ends with Linear issues carrying acceptance criteria,
plus an ADR if the planning surfaced a real decision. It does not end with a
step list in a transcript that scrolls away.

### Ask for decisions, not for a walkthrough

The model sequences work well and chooses between options poorly. You are the
reverse. Point planning at the gap: ask which choices this work forces and
what each one costs. Three trade-offs is ten lines you will actually read.
Forty steps is a page you will skim and wave through.

### How the two sessions run

- **Planning.** Use plan mode (`Shift+Tab`, or `/plan`). Claude explores and
  proposes but cannot edit until you approve. Output: issues with criteria,
  and any ADR the planning turned up.
- **Building.** A clean session, one issue, per the boundary rule above.

Plan length is a signal, not a nuisance. A plan that outgrows one screen means
the issue is too big. Split the issue, not the plan.

A decision the planning surfaced has to land as an ADR before building starts.
Left inside the plan, it lives only in a transcript, and the concept drifts
even while the code stays fine.

### What earns an ADR

An ADR exists to stop a decision being quietly re-litigated by someone who no
longer holds the reason — usually you, months later, or a model starting from
zero. So the test is what could undo it by accident.

A choice contained in one file has nothing to protect: changing it means opening
the file that states it, and the comment there is already the argument. A choice
that constrains code not yet written is the opposite. Someone violates it who
never saw it, in a file that didn't exist when it was made. That one needs an
entry.

The failure mode in the other direction is granularity. Recording every small
call produces a decision log nobody reads, which protects nothing. Decisions
taken together about one thing are one entry.

### Why a rule is a separate thing from an ADR

An ADR answers "why is it like this." A rule answers "what do I do here." They
get separated because they are needed at different moments, by readers in
different states.

You need the reasoning when you are considering a change to the decision —
rarely, and you go looking for it. You need the constraint every time you touch
a governed file — often, without knowing to ask. One is pulled, the other has to
be pushed. `paths:` frontmatter is the push: a rule loads when a matching file is
opened, so the constraint arrives without anyone remembering it exists.

That is why a decision that applies more than once creates a rule, and why they
cite each other. A rule with no source is a constraint whose reason has been
lost, and the first person to find it inconvenient will delete it. A decision
with no rule, where one is warranted, will be violated by someone who never read
it.

The source does not have to be an ADR. A product constraint in `01-product.md`
creates a rule the same way. `.claude/rules/thechat-domain.md` carries one of
each.

### What rules actually buy, and what they don't

Two limits, stated plainly, because a mechanism you overestimate is worse than
one you know is partial.

**A rule is context, not enforcement.** Claude Code's own documentation is
explicit: instructions shape behavior but are not an enforcement layer. To block
an action regardless of what the model decides, you need a `PreToolUse` hook.
A rule raises the odds; it does not close the door.

**A path-scoped rule loads on read, not on write.** It fires when a matching file
is opened, not when a new matching file is created. So a rule governing
`internal/thechat/**/*.go` reaches you while editing an existing entity, but not
necessarily while adding the first file for a new one. In practice the pattern
gets read before it gets copied, which usually covers it — but "usually" is the
honest word.

`.claude/rules/adr.md` is scoped to `docs/02-system.md` for exactly this reason:
ADRs live in that file, so writing one means opening it, and opening it loads the
template. The trigger and the work are the same action.

Both limits point the same way. Rules make the constraint likely to arrive, not
certain. Certainty needs a hook or a CI check, which brings us to what's missing.

---

## What to plan at

### The unit is the unit of integration

"Add a user model," "add a test," and "add a feature" are not three sizes of
one thing. The first is a domain change, the second is part of some other
unit's definition of done, and the third is usually several issues.

The axis that matters: integration is the part that has to stay ordered, so
plan in units that can merge on their own and leave main working. If it
cannot ship alone, it is not a unit.

The working test is whether you can write acceptance criteria **before** you
start digging. If you can't, the unit is too big — it needs exploration
first, and exploration is already building.

### Slice vertically, never by layer

Splitting by layer (model → repository → handler → UI) produces issues that
can't be verified one at a time. Their criteria collapse into "the code
exists," which is not a truth condition. Everything only proves out once all
of them land, so drift is guaranteed by construction.

A vertical slice is the smallest change that runs end to end and can be
watched working. It touches storage, API, and UI at once, and it has one
observable criterion. Vertical criteria hold regardless of state. Layer
criteria do not.