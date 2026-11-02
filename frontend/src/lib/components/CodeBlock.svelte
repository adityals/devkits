<script lang="ts" module>
    import { createHighlighterCore } from "shiki/core";
    import { createJavaScriptRegexEngine } from "shiki/engine/javascript";

    export type Lang = "css" | "html" | "javascript" | "json" | "xml" | "text";

    // Highlighting runs on the UI thread, so large bodies render as plain text.
    const MAX_HIGHLIGHT = 100_000;

    const highlighter = createHighlighterCore({
        langs: [
            import("shiki/langs/css.mjs"),
            import("shiki/langs/html.mjs"),
            import("shiki/langs/javascript.mjs"),
            import("shiki/langs/json.mjs"),
            import("shiki/langs/xml.mjs"),
        ],
        themes: [import("shiki/themes/github-dark.mjs")],
        engine: createJavaScriptRegexEngine(),
    });
</script>

<script lang="ts">
    type Props = {
        code: string;
        lang?: Lang;
        class?: string;
    };

    const {
        code,
        lang = "text",
        class: className = "",
    }: Props = $props();
</script>

<div
    class={[
        "rounded-container overflow-auto text-xs [&>pre]:p-3 [&>pre]:w-fit [&>pre]:min-w-full [&>pre]:min-h-full",
        className,
    ]}
>
    {#if code.length > MAX_HIGHLIGHT}
        <pre class="bg-surface-950 text-surface-50">{code}</pre>
    {:else}
        {#await highlighter then h}
            <!-- Shiki escapes the code, so untrusted bodies are safe here. -->
            {@html h.codeToHtml(code, {
                lang,
                theme: "github-dark",
            })}
        {/await}
    {/if}
</div>
