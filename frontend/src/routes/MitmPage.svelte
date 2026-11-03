<script lang="ts">
    import type { Component } from "svelte";
    import { Tabs } from "@skeletonlabs/skeleton-svelte";
    import Control from "../modules/mitm/proxy/Control.svelte";
    import Traffic from "../modules/mitm/traffic/Traffic.svelte";

    type LazyTab = {
        label: string;
        load: () => Promise<{ default: Component }>;
    };

    const lazyTabs: Record<string, LazyTab> = {
        rules: {
            label: "Rule Engine",
            load: () => import("../modules/mitm/rule/Rules.svelte"),
        },
    };

    const visited = $state<Record<string, boolean>>({});
</script>

<div class="flex h-full flex-col">
    <Control />
    <main class="min-h-0 flex-1 p-4">
        <Tabs
            defaultValue="traffic"
            onValueChange={(details) => {
                visited[details.value] = true;
            }}
            class="flex h-full flex-col"
        >
            <Tabs.List>
                <Tabs.Trigger value="traffic">Network Traffic</Tabs.Trigger>
                {#each Object.entries(lazyTabs) as [value, tab] (value)}
                    <Tabs.Trigger {value}>{tab.label}</Tabs.Trigger>
                {/each}
                <Tabs.Indicator />
            </Tabs.List>
            <Tabs.Content value="traffic" class="min-h-0 flex-1">
                <Traffic />
            </Tabs.Content>
            {#each Object.entries(lazyTabs) as [value, tab] (value)}
                <Tabs.Content {value} class="min-h-0 flex-1">
                    {#if visited[value]}
                        {#await tab.load() then { default: Panel }}
                            <Panel />
                        {:catch error}
                            <p class="card preset-tonal-error p-3 text-sm">
                                {String(error)}
                            </p>
                        {/await}
                    {/if}
                </Tabs.Content>
            {/each}
        </Tabs>
    </main>
</div>
