<script lang="ts">
    import { LoaderCircle } from "@lucide/svelte";

    const { class: className = "" }: { class?: string } = $props();

    const DELAY_MS = 150;

    let visible = $state(false);

    $effect(() => {
        const timer = setTimeout(() => (visible = true), DELAY_MS);

        return () => clearTimeout(timer);
    });
</script>

{#if visible}
    <div
        class={[
            "text-surface-600-400 flex items-center gap-2 text-sm",
            className,
        ]}
        role="status"
    >
        <LoaderCircle size={16} class="animate-spin" />
    </div>
{/if}
