---
paths:
  - "internal/thechat/**/*.go"
---

# Domain constraints

Two live constraints on `package thechat`. Read the source before working around
either: one is a system decision, the other is a product constraint.

## Entity identity — from ADR 0001

Every entity ID is an alias of the shared generic wrapper:

```go
type thing struct{}          // phantom type, unexported, empty
type ThingID = ID[thing]     // alias, not a defined type

func NewThingID() ThingID {
    return ID[thing]{uuid.Must(uuid.NewV7())}
}
```

- UUIDv7, never v4 and never a database integer. It is time-ordered, so it
  doesn't fragment the index and carries creation order without a column.
- `uuid.Must`. `NewV7` fails only if the OS CSPRNG is broken, which is
  unrecoverable, so it fails fast rather than threading an error through every
  call site.
- Alias (`=`), not a defined type. A defined type would need its own methods to
  reach `uuid.UUID`.
- Never add a hand-written per-entity ID type. One generic definition covers all
  of them.

## No container entity — from `01-product.md` §Constraints

Never add a `Chat`, `Conversation`, `Thread`, or `Session` entity. Messages
reference `UserID` directly.

Adding one does not extend this design, it reverses the product constraint. That
is a product call, so it changes `01-product.md` first.
