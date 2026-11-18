<script lang="ts">
    import { ArrowRight, Copy } from "@lucide/svelte";
    import "prism-code-editor/prism/languages/json";
    import CodeBlock, {
        type Lang as BlockLang,
    } from "../../lib/components/CodeBlock.svelte";
    import CodeEditor from "../../lib/components/CodeEditor.svelte";
    import { toaster } from "../../lib/toaster";
    import { Converter } from "./converter.svelte";
    import { sources } from "./provider/registry";
    import type { Lang } from "./provider/types";

    const labels: Record<Lang, string> = {
        json: "JSON",
        go: "Go",
        ts: "TypeScript",
    };

    const blockLangs: Record<Lang, BlockLang> = {
        json: "json",
        go: "go",
        ts: "typescript",
    };

    const converter = new Converter();
    let pasted = false;

    const handleChange = (next: string) => {
        converter.input = next;
        if (pasted) {
            pasted = false;
            converter.run();
        }
    };

    const handleKeydown = (event: KeyboardEvent) => {
        if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
            event.preventDefault();
            event.stopPropagation();
            converter.run();
        }
    };

    const handleFromChange = () => {
        converter.to = converter.targets[0];
        converter.run();
    };

    const handleCopy = async () => {
        await navigator.clipboard.writeText(converter.output);
        toaster.success({ title: "Copied" });
    };
</script>

<div
    class="flex h-full flex-col gap-3"
    onkeydowncapture={handleKeydown}
    role="presentation"
>
    <div class="flex items-center gap-2">
        <select
            class="select w-auto"
            aria-label="From"
            bind:value={converter.from}
            onchange={handleFromChange}
            disabled={sources.length < 2}
        >
            {#each sources as lang (lang)}
                <option value={lang}>{labels[lang]}</option>
            {/each}
        </select>
        <ArrowRight size={16} />
        <select
            class="select w-auto"
            aria-label="To"
            bind:value={converter.to}
            onchange={() => converter.run()}
        >
            {#each converter.targets as lang (lang)}
                <option value={lang}>{labels[lang]}</option>
            {/each}
        </select>
        <input
            class="input w-40"
            aria-label="Type name"
            placeholder="Type name"
            autocapitalize="none"
            bind:value={converter.typeName}
        />
        {#if converter.to === "go"}
            <input
                class="input w-32"
                aria-label="Go package"
                placeholder="Package"
                bind:value={converter.goPackage}
            />
        {/if}
        <button
            type="button"
            class="btn preset-filled-primary-500"
            title="Convert (⌘/Ctrl+Enter)"
            onclick={() => converter.run()}
        >
            Convert
        </button>
        <button
            type="button"
            class="btn-icon preset-tonal ml-auto"
            aria-label="Copy output"
            title="Copy output"
            disabled={!converter.output || Boolean(converter.error)}
            onclick={handleCopy}
        >
            <Copy size={16} />
        </button>
    </div>
    <div class="grid min-h-0 flex-1 grid-cols-2 gap-3">
        <div
            class="min-h-0"
            onpaste={() => (pasted = true)}
            role="presentation"
        >
            <CodeEditor
                value={converter.input}
                language={converter.from}
                onChange={handleChange}
            />
        </div>
        {#if converter.error}
            <p class="card preset-tonal-error h-fit p-3 text-sm">
                {converter.error}
            </p>
        {:else if converter.output}
            <CodeBlock
                class="h-full"
                code={converter.output}
                lang={blockLangs[converter.to]}
            />
        {:else}
            <p class="text-surface-600-400 self-center text-center text-sm">
                Paste {labels[converter.from]} to generate types
            </p>
        {/if}
    </div>
</div>
