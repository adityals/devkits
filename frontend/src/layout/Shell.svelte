<script lang="ts">
    import type { Snippet } from "svelte";
    import { AppBar } from "@skeletonlabs/skeleton-svelte";
    import { Wrench } from "@lucide/svelte";

    const { children }: { children: Snippet } = $props();
</script>

<div class="grid h-screen grid-rows-[auto_1fr]">
    <AppBar class="border-surface-200-800 border-b">
        <AppBar.Toolbar class="grid-cols-[auto_1fr]">
            <AppBar.Lead>
                <a
                    href="#/"
                    class="flex items-center gap-2"
                >
                    <Wrench size={20} />
                    <span class="h6">DevKits</span>
                </a>
            </AppBar.Lead>
        </AppBar.Toolbar>
    </AppBar>
    <div class="min-h-0 min-w-0 overflow-auto">
        {@render children()}
    </div>
</div>
