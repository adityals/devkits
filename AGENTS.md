# AGENTS.md

DevKits: Wails v3 desktop app. Go backend (`main.go`, `internal/`), Svelte 5 + TypeScript frontend (`frontend/`, pnpm, Vite, Tailwind, Skeleton).

## Directory Layout

The app is split into modules (e.g. `mitm`, `imgcompress`, `jsonstruct`). Use the same module name on both sides.

```
main.go                         # wiring only: create app, register each module's Services
internal/
  <module>/
    <module>.go                 # Services(deps...) []application.Service + RegisterEvent in init
    <sub>/                      # one package per concern, e.g. mitm/{proxy,rule,traffic,cert}
frontend/src/
  App.svelte                    # Shell + Router
  routes.ts                     # one route per module: /<module>
  routes/<Module>Page.svelte    # one page per module, plus NotFound.svelte
  layout/Shell.svelte           # navigation rail; add one link per module
  lib/components/               # UI shared across modules
  modules/<module>/<sub>/       # module components and .svelte.ts logic, mirrors internal/<module>/<sub>
```

- Modules never import each other. Shared code goes in `lib/` (frontend) or a specifically named `internal/` package once two modules need it; never `utils`/`common`.
- Prefix event names with the module: `<module>:<sub>:<event>` (e.g. `mitm:traffic:flow`).
- Moving a Go package changes binding import paths; regenerate bindings and fix imports in the same change.

## Code Write Practices

### Go (`*.go`)
- Follow Go best practices and idiomatic Go (Effective Go, Go Code Review Comments).
- Use dependency injection: pass dependencies through constructors (`NewXxx(deps...)`), depend on small interfaces, no package-level mutable state or hidden globals.
- Wrap errors with context (`fmt.Errorf("...: %w", err)`); never ignore errors.

### Wails v3
- Follow Wails v3 best practices.
- Expose backend logic as services registered with the application; keep services thin and delegate to `internal/` packages.
- Use typed events for backend-to-frontend communication.
- Never edit generated `frontend/bindings/`; regenerate them after changing service signatures.

### TypeScript (`*.ts`)
- Follow TypeScript best practices and idiomatic TypeScript.
- Strict typing; no `any`. Prefer generated binding types over hand-written duplicates.

### Svelte (`*.svelte`)
- Follow Svelte 5 best practices and idiomatic Svelte: runes (`$state`, `$derived`, `$effect`, `$props`), `<script lang="ts">`.
- Keep components small; move shared logic into `.ts` modules.

### UI Design System (Skeleton)
- All frontend UI follows Skeleton for Svelte. Reference: https://www.skeleton.dev/llms.txt (Svelte docs under `/docs/svelte/`).
- Use `@skeletonlabs/skeleton-svelte` components before hand-rolling one (`import { Avatar } from '@skeletonlabs/skeleton-svelte'`, composed parts like `<Avatar.Image />`).
- Style with Skeleton Tailwind utilities and presets (`btn`, `card`, `input`, `preset-filled`, `preset-tonal`, `preset-outlined`, `preset-filled-primary-500`), not custom CSS.
- Use theme color tokens (`primary`, `secondary`, `tertiary`, `success`, `warning`, `error`, `surface`) and color pairings (`bg-surface-50-950`, `text-surface-950-50`); never hard-coded hex or raw Tailwind palette colors.
- Theme is `cerberus`, set with `data-theme` on `<html>` (`frontend/index.html`) and imported in `frontend/src/app.css`. Dark mode uses Tailwind's `dark:` variant.

### Formatting (Go, TS, Svelte, JS)
- Use line breaks generously: one element per line in multi-element slice/map/array/object literals, function arguments, and parameter lists, with trailing commas where the language allows.
- Separate logical blocks with a blank line:
  - after a closing `}` of an `if`/`for`/`switch` when more statements follow
  - between a group of assignments and the `if` that checks them, unless it is a single `x, err := f()` + `if err != nil` pair, which stays together
  - before the final `return` when it follows a block
- Keep tightly bound lines together: `x, err := f()` with its `if err != nil`, and `defer cleanup()` right after the acquire's error check.

### Tests
- Only create test files for critical business paths.

### Comments
- Do not write unnecessary comments. Comment only the why the code cannot show.

## Work With Human Rules
- Review budget is 150 LOC total across all files (tests included) since the last human review. Before a write would push the running total past 150, stop and ask the human to review. Continue only after approval, then reset the count.
- The budget is only a pause point. Never compress, inline, or reformat code to fit it. If the next piece doesn't fit, stop and put it in the next chunk; one file may span several chunks.
