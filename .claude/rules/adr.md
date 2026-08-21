---
paths:
  - "docs/02-system.md"
---

# Writing an ADR

ADRs are numbered `###` entries under `## Decisions` in `docs/02-system.md`.
Not separate files. Copy the shape below.

```markdown
### NNNN — Title as a claim, present tense

- Status: Accepted, YYYY-MM-DD
- Job: job from `01-product.md`, or `none (structural)`
- Terms: domain terms this touches, from that file's Domain model section
- Rule: `.claude/rules/<file>.md`, or `none`

**Context.** What forced the decision. Which constraint or conflict brought it up.

**Decision.** What we chose. One paragraph, present tense.

**Rejected.** What else was on the table and why it lost. Without this the choice
reads as arbitrary six months from now.

**Costs.** The price we accepted, stated plainly. Never "no downsides" — if you
can't see one, the decision hasn't been thought through yet.

**Revisit when.** Concrete, observable conditions. Not "if things change." Write
"when we move to a database with no native uuid type," or "if ID generation shows
up in a profile." This is the most important field: it turns the decision from a
rule you obey into a claim that can be proven wrong.
```

## Conventions

**Status moves one way only:** `Accepted` → `Superseded by NNNN`.

**Never edit the decision text of an accepted ADR.** A changed decision is a new
entry with the next number. When it supersedes an older one, edit exactly one
line in the old entry — its Status — and nothing else. The old reasoning stays
readable, because the point of the record is what someone believed at the time.

**If the decision applies more than once, it needs a rule.** Write it under
`.claude/rules/`, path-scoped to the files it governs, then link both ways: the
ADR's `Rule:` field points at the rule, and the rule cites this ADR's number. An
ADR is the decision; a rule is the live constraint that decision creates.

A rule's source is not always an ADR. A product constraint in `01-product.md`
§Constraints creates rules the same way, and the rule cites whichever it came
from.

**A record filed in the wrong place is deleted, not superseded.** Superseding is
for a decision that was right and stopped being right. An entry that was never a
system decision — a product call, say — just goes, and its number is free again.

**A rule pointing at a deleted or superseded source is a bug.** Nothing checks
this yet — see `00-process.md` §4.

## What does not go here

A choice contained in one file is a comment, not a decision. If reversing it
means editing the file that states it, the comment there is already the record.
Decisions taken together about one thing are one entry, not several.
