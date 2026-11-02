<script lang="ts">
    import { Tabs } from "@skeletonlabs/skeleton-svelte";
    import Control from "../modules/mitm/proxy/Control.svelte";
    import Traffic from "../modules/mitm/traffic/Traffic.svelte";
</script>

<div class="flex h-full flex-col">
    <Control />
    <main class="min-h-0 flex-1 p-4">
        <Tabs defaultValue="traffic" class="flex h-full flex-col">
            <Tabs.List>
                <Tabs.Trigger value="traffic">Network Traffic</Tabs.Trigger>
                <Tabs.Trigger value="rules">Rule Engine</Tabs.Trigger>
                <Tabs.Indicator />
            </Tabs.List>
            <Tabs.Content value="traffic" class="min-h-0 flex-1">
                <Traffic />
            </Tabs.Content>
            <Tabs.Content value="rules" class="min-h-0 flex-1"></Tabs.Content>
        </Tabs>
    </main>
</div>
