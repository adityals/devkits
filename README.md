# DevKits

A desktop toolkit for developers, built with [Wails v3](https://v3.wails.io/) (Go backend) and Svelte 5 (frontend).

## Modules

- [MITM Proxy](docs/mitm.md): record HTTP(S) traffic and rewrite it with JS rules.
- [Type Generator](docs/typegen.md): generate typed code from data or schemas.

## Development

Requirements: Go 1.26, Node.js with pnpm, and the `wails3` CLI.

```sh
wails3 task dev      # run with hot reload
wails3 task build    # build for the host OS
wails3 task test     # run Go tests in ./internal/...
```

Pass extra `go test` arguments after `--`, for example `wails3 task test -- -run TestFlow ./internal/mitm/rule`.

## Layout

```
main.go                  # app wiring
internal/<module>/       # Go services, one package per concern
frontend/src/routes/     # one page per module
frontend/src/modules/    # module components and logic
```

See `AGENTS.md` for conventions.
