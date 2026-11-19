# Type Generator

Generates typed code from data or schemas. Alpha routes:

- JSON to Go and JSON to TypeScript, powered by [quicktype](https://github.com/glideapps/quicktype) in the frontend.
- Go to TypeScript, powered by [tygo](https://github.com/gzuidhof/tygo) in the Go backend. Paste type declarations; a `package` line is optional. `json` tags set field names, and pointers and `omitempty` become optional fields.

- Paste input to convert right away. After typing, convert with **Convert** or `⌘/Ctrl+Enter`; typing alone never converts, so half-written input doesn't produce errors.
- For JSON input, set the top-level type name (default `Root`) and, for Go, the package name (default `main`).
- Output is types only: Go structs with `json` tags, or TypeScript interfaces, with no marshalling helpers.
- `quicktype-core` (~4 MB) loads on the first conversion, so the first run is slower.

## Adding a conversion

Each conversion is a route in `frontend/src/modules/typegen/provider/registry.ts`, keyed `<from>_<to>` and mapped to a provider:

```ts
export const routes: Partial<Record<Route, Provider>> = {
  json_go: quicktype,
  json_ts: quicktype,
  go_ts: backend,
};
```

A provider is a class that implements `Provider.convert(input, from, to, opts)`. The From and To dropdowns are built from the route keys, and a missing route fails with `no provider for <from> → <to>`. Conversions that need Go use `BackendProvider`, which calls the `convert` Wails service. On the Go side, add a `convert.Converter` to the route map in `internal/typegen/typegen.go`.
