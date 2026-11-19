# MITM Proxy

An HTTP(S) proxy that records traffic and rewrites it with JavaScript rules.

- **Proxy** listens on port `39080` on all interfaces, so other devices on your network can use it.
- **CA certificate** is created in `~/.config/devkits/ca/` on first use. Trust it from the app to inspect HTTPS. On macOS it goes into the login keychain.
- **System proxy** can point macOS at `127.0.0.1:39080`. It is turned off again when the proxy stops or the app quits.
- **Traffic** keeps the latest 500 flows in memory. Bodies over 1 MB are truncated; bodies over 5 MB are streamed, so the flow is recorded without its body.

Trusting the CA and the system proxy are macOS only for now.

Because the proxy listens on all interfaces, anyone on the same network can use it while it runs. Stop it on networks you don't trust.

## Proxying an Android device

1. Connect the device to the same network as your computer.
2. In the Wi-Fi settings, set a manual proxy: your computer's LAN IP and port `39080`.
3. Copy `~/.config/devkits/ca/mitmproxy-ca-cert.pem` to the device and install it under **Settings → Security → Encryption & credentials → Install a certificate → CA certificate**.

On Android 7 and later, apps ignore user-installed CAs by default, so their HTTPS traffic is not intercepted. Chrome trusts them. For your own app, trust user CAs in a debug-only `network_security_config`:

```xml
<network-security-config>
  <debug-overrides>
    <trust-anchors>
      <certificates src="user" />
    </trust-anchors>
  </debug-overrides>
</network-security-config>
```

## Rules

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
