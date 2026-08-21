# Product and jobs to be done

## Product overview

The Chat is one limitless conversation with LLM models.

There is no "New Chat" button and no list of threads. There is one conversation,
and it keeps going. What the interface shows instead is how much of the context
window is in use, and a **Clear** button that empties it.

The bet is that "start a new chat" is the wrong control. It asks you to decide,
up front, whether what you are about to say belongs with what came before. That
is a decision you usually can't make yet, and it leaves you with a list of
half-finished threads either way. Showing context usage replaces that guess with
a reading, and Clear turns it into an act you take when you can see the reason.

This is principles 1 and 3 applied to the product surface. The model's state is
bounded, so make the bound visible instead of leaving it to be inferred from
degrading answers.

## Jobs To Be Done

- **Keep going without deciding where this belongs.** Say the next thing without
  first judging whether it's the same topic, a new one, or worth a new thread.
- **Know how much room is left.** See context usage before quality starts
  degrading, not after — while there's still time to act on it.
- **Drop the context on purpose.** Clear when the subject has actually changed,
  as a deliberate act with a visible reason, not as a side effect of starting
  something new.

## Constraints

Things the product does not do, and won't. Each is a choice, not a gap.

**One conversation per user. No threads, no sessions, no archive.** The bet above
is what this rests on.

*What it costs.* No history after a Clear, so the button is destructive and has
to feel that way. No way to keep two subjects alive at once — people who want
that will ask for threads. Search across past conversations is off the table by
construction.

*What would change our mind.* Clear being avoided rather than used. That means
people are hoarding context they can't afford to lose, and the missing feature is
retention, not threads.

*What the domain does about it.* No container entity — messages hang off the
user. Carried by `.claude/rules/thechat-domain.md`.
