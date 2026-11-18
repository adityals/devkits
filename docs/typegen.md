# Type Generator

Generates typed code from data or schemas. Alpha: JSON to Go and JSON to TypeScript, powered by [quicktype](https://github.com/glideapps/quicktype) and run entirely in the frontend.

- Paste input to convert right away. After typing, convert with **Convert** or `⌘/Ctrl+Enter`; typing alone never converts, so half-written input doesn't produce errors.
- Set the top-level type name (default `Root`) and, for Go, the package name (default `main`).
- Output is types only: Go structs with `json` tags, or TypeScript interfaces, with no marshalling helpers.
- `quicktype-core` (~4 MB) loads on the first conversion, so the first run is slower.

## Adding a conversion

Each conversion is a route in `frontend/src/modules/typegen/provider/registry.ts`, keyed `<from>_<to>` and mapped to a provider:

```ts
export const routes: Partial<Record<Route, Provider>> = {
  json_go: quicktype,
  json_ts: quicktype,
};
```

A provider is a class that implements `Provider.convert(input, from, to, opts)`. The From and To dropdowns are built from the route keys, and a missing route fails with `no provider for <from> → <to>`. Conversions that need Go (such as Go to TypeScript or protobuf) can be served by a provider that calls a Wails service.
