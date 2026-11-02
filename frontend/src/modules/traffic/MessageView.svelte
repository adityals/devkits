<script lang="ts">
    import CodeBlock from "../../lib/components/CodeBlock.svelte";
    import type { Message } from "../../../bindings/devkits/internal/traffic";
    import { formatBody, headerValue, langFor } from "./content-type";

    const { message }: { message: Message } = $props();

    const lang = $derived(langFor(headerValue(message.header, "Content-Type")));
    const headers = $derived(
        Object.entries(message.header ?? {}).toSorted(([a], [b]) =>
            a.localeCompare(b),
        ),
    );
</script>

<div class="space-y-4">
    <table class="table text-xs">
        <tbody>
            {#each headers as [name, values] (name)}
                <tr>
                    <th class="w-1/3 align-top font-semibold">{name}</th>
                    <td class="break-all">{values?.join(", ")}</td>
                </tr>
            {/each}
        </tbody>
    </table>
    {#if message.streamed}
        <p class="text-surface-600-400 text-sm">Body streamed, not captured.</p>
    {:else if message.body}
        {#if message.truncated}
            <p class="text-warning-500 text-sm">Body truncated.</p>
        {/if}
        <!-- 10 visible lines plus the p-3 padding, scroll for the rest. -->
        <CodeBlock
            code={formatBody(message.body, lang)}
            {lang}
            class="h-[calc(10lh+1.5rem)]"
        />
    {/if}
</div>
