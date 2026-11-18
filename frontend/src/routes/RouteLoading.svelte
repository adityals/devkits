<script lang="ts">
    import { LoaderCircle } from "@lucide/svelte";

    // Embedded chunks usually load in a frame, so wait before showing anything to avoid a flash.
    const DELAY_MS = 150;

    let visible = $state(false);

    $effect(() => {
        const timer = setTimeout(() => (visible = true), DELAY_MS);

        return () => clearTimeout(timer);
    });
</script>

{#if visible}
    <div
        class="text-surface-600-400 flex h-full items-center justify-center gap-2 text-sm"
        role="status"
    >
        <LoaderCircle
            size={16}
            class="animate-spin"
        />
        Loading…
    </div>
{/if}
