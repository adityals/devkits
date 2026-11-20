<script lang="ts" module>
    import { createHighlighterCore } from "shiki/core";
    import { createJavaScriptRegexEngine } from "shiki/engine/javascript";

    export type Lang =
        | "css"
        | "go"
        | "html"
        | "javascript"
        | "json"
        | "typescript"
        | "xml"
        | "text";

    // Highlighting runs on the UI thread, so large bodies render as plain text.
    const MAX_HIGHLIGHT = 100_000;

    // shortcut: theme is picked once from the OS setting like CodeEditor, add a listener if live switching matters.
    const theme = window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "github-dark"
        : "github-light";

    const highlighter = createHighlighterCore({
        langs: [
            import("shiki/langs/css.mjs"),
            import("shiki/langs/go.mjs"),
            import("shiki/langs/html.mjs"),
            import("shiki/langs/javascript.mjs"),
            import("shiki/langs/json.mjs"),
            import("shiki/langs/typescript.mjs"),
            import("shiki/langs/xml.mjs"),
        ],
        themes: [
            theme === "github-dark"
                ? import("shiki/themes/github-dark.mjs")
                : import("shiki/themes/github-light.mjs"),
        ],
        engine: createJavaScriptRegexEngine(),
    });
</script>

<script lang="ts">
    type Props = {
        code: string;
        lang?: Lang;
        class?: string;
    };

    const { code, lang = "text", class: className = "" }: Props = $props();
</script>

<div
    class={[
        "rounded-container overflow-auto text-xs [&>pre]:p-3 [&>pre]:w-fit [&>pre]:min-w-full [&>pre]:min-h-full",
        className,
    ]}
>
    {#if code.length > MAX_HIGHLIGHT}
        <pre class="bg-surface-50-950 text-surface-950-50">{code}</pre>
    {:else}
        {#await highlighter then h}
            <!-- Shiki escapes the code, so untrusted bodies are safe here. -->
            {@html h.codeToHtml(code, {
                lang,
                theme,
            })}
        {/await}
    {/if}
</div>
