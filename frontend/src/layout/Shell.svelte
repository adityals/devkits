<script lang="ts">
    import type { Snippet } from "svelte";
    import { Navigation } from "@skeletonlabs/skeleton-svelte";
    import { Network } from "@lucide/svelte";
    import { router } from "svelte-spa-router";

    const { children }: { children: Snippet } = $props();

    const links = [
        {
            href: "/mitm",
            label: "MITM",
            icon: Network,
        },
    ];
</script>

<div class="grid h-screen grid-cols-[auto_1fr]">
    <Navigation
        layout="rail"
        class="border-surface-200-800 border-r"
    >
        <Navigation.Content>
            <Navigation.Menu>
                {#each links as link (link.href)}
                    {@const active = router.location.startsWith(link.href)}
                    <Navigation.TriggerAnchor
                        href={`#${link.href}`}
                        aria-current={active ? "page" : undefined}
                        class={[active && "preset-tonal-primary"]}
                    >
                        <link.icon size={20} />
                        <Navigation.TriggerText>{link.label}</Navigation.TriggerText>
                    </Navigation.TriggerAnchor>
                {/each}
            </Navigation.Menu>
        </Navigation.Content>
    </Navigation>
    <div class="min-h-0 min-w-0">
        {@render children()}
    </div>
</div>
