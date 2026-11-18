# DevKits

A desktop toolkit for developers, built with [Wails v3](https://v3.wails.io/) (Go backend) and Svelte 5 (frontend).

## Modules

### MITM

An HTTP(S) proxy that records traffic and rewrites it with JavaScript rules.

- **Proxy** listens on port `8000` on all interfaces, so other devices on your network can use it.
- **CA certificate** is created in `~/.config/devkits/ca/` on first use. Trust it from the app to inspect HTTPS. On macOS it goes into the login keychain.
- **System proxy** can point macOS at `127.0.0.1:8000`. It is turned off again when the proxy stops or the app quits.
- **Traffic** keeps the latest 500 flows in memory. Bodies over 1 MB are truncated; bodies over 5 MB are streamed, so the flow is recorded without its body.

Trusting the CA and the system proxy are macOS only for now.

#### Rules

Each rule is a JS file in `~/.config/devkits/proxy_rules/`: `<id>.js` is enabled, `_<id>.js` is disabled. Up to 5 rules can be enabled at once.

```js
const rule = {
  host: ["api.example.com", "*.example.dev"],
  modifier: "^~",
  match: "/v1/",
  onRequestHeader(flow, h) {
    h.set("X-Debug", "1");
  },
  onResponseBody(flow, body) {
    return body.replace("prod", "local");
  },
};
```

- `host` (optional): a string or an array. Accepts `example.com`, `example.com:8443`, `*.example.com` and `www.example.*`. Omit it to match any host.
- `match` and `modifier` pick the path, with nginx `location` priority:

  | `modifier` | Meaning |
  | --- | --- |
  | `=` | exact path |
  | `^~` | prefix; skips regex rules when it is the longest prefix |
  | `~` | regex |
  | `~*` | case-insensitive regex |
  | omitted | prefix |

- Hooks: define at least one. All of them receive `flow` first.

  | Hook | Second argument | Return |
  | --- | --- | --- |
  | `onRequestHeader` | header object | ignored |
  | `onResponseHeader` | header object | ignored |
  | `onRequestBody` | body string | new body, or `undefined` to keep it |
  | `onResponseBody` | body string | new body, or `undefined` to keep it |

- `flow` has `url`, `host`, `path` and `query` (an object of string arrays, e.g. `flow.query.q[0]`).
- The header object has `get`, `values`, `set`, `add`, `del` and `toObject`.
- `console.log` and the other `console` methods write to the app log.
- Each hook call times out after 500 ms.

A body hook makes the proxy buffer and decode every matched body, so only define the hooks you use.

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
