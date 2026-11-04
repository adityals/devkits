<script lang="ts">
    import { Network } from "@lucide/svelte";

    const modules = [
        {
            href: "/mitm",
            label: "MITM Proxy",
            description: "Record HTTP traffic and rewrite it with JS rules.",
            icon: Network,
        },
    ];
</script>

<div class="mx-auto max-w-5xl space-y-8 p-6">
    <header class="space-y-2 py-8 text-center">
        <h1 class="h2">DevKits</h1>
        <p class="text-surface-600-400">Developer tools in one desktop app.</p>
    </header>
    <nav class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {#each modules as module (module.href)}
            <a
                href={`#${module.href}`}
                class="card preset-outlined-surface-200-800 hover:preset-tonal-primary space-y-2 p-4 transition-colors"
            >
                <module.icon size={24} />
                <h2 class="h5">{module.label}</h2>
                <p class="text-surface-600-400 text-sm">{module.description}</p>
            </a>
        {/each}
    </nav>
</div>
