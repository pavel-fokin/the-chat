# The Chat

**The Chat** is one limitless conversation with LLM models — no "New Chat"
button, no thread list. The interface shows how much of the context window is in
use, and a Clear button empties it. See `docs/01-product.md`.

## Run it

Requires Go at the version in `go.mod`, and Node at the current LTS. Both track
upstream rather than a number pinned here, so this line doesn't need editing when
a new LTS lands. The build's real floor comes from Vite (`^20.19.0 || >=22.12.0`)
and every LTS since has cleared it.

```sh
make install   # npm install in web/
make build     # builds the frontend, then the Go binary
make run       # build, then serve on 0.0.0.0:8080
```

Then open http://localhost:8080. Override the address with
`./bin/the-chat-server -addr 127.0.0.1:3000`.

**The frontend has to be built before the Go binary.** `go:embed` reads
`web/dist` at compile time, so `go build` on its own will use a stale `dist` — or
fail if there isn't one. `make build` orders them correctly; use it.

For frontend-only work, `npm --prefix web run dev` gives you Vite with hot
reload, but it does not serve the Go backend.
