# AGENTS.md

DevKits: Wails v3 desktop app. Go backend (`main.go`, `internal/`), Svelte 5 + TypeScript frontend (`frontend/`, pnpm, Vite, Tailwind, Skeleton).

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
