<script lang="ts">
    import { Tabs } from "@skeletonlabs/skeleton-svelte";
    import { X } from "@lucide/svelte";
    import { Get } from "../../../bindings/devkits/internal/traffic/trafficservice";
    import type { Flow } from "../../../bindings/devkits/internal/traffic";
    import MessageView from "./MessageView.svelte";

    type Props = {
        id: string;
        onclose: () => void;
    };

    const { id, onclose }: Props = $props();

    let flow = $state<Flow | null>(null);
    let error = $state("");

    $effect(() => {
        let stale = false;

        Get(id)
            .then((f) => {
                if (!stale) {
                    flow = f;
                    error = "";
                }
            })
            .catch((e) => {
                if (!stale) {
                    error = String(e);
                }
            });

        return () => {
            stale = true;
        };
    });
</script>

<aside class="card preset-outlined-surface-200-800 flex min-h-0 flex-col">
    <header class="border-surface-200-800 flex items-center gap-2 border-b p-3">
        <p
            class="min-w-0 flex-1 truncate text-sm font-semibold"
            title={flow?.summary.url}
        >
            {flow?.summary.method}
            {flow?.summary.url}
        </p>
        <button
            type="button"
            class="btn-icon btn-icon-sm hover:preset-tonal"
            aria-label="Close"
            title="Close"
            onclick={onclose}
        >
            <X size={16} />
        </button>
    </header>
    {#if error}
        <p class="preset-tonal-error p-3 text-sm">{error}</p>
    {:else if flow}
        <Tabs
            defaultValue="response"
            class="flex min-h-0 flex-1 flex-col p-3"
        >
            <Tabs.List>
                <Tabs.Trigger value="request">Request</Tabs.Trigger>
                <Tabs.Trigger value="response">Response</Tabs.Trigger>
                <Tabs.Indicator />
            </Tabs.List>
            <Tabs.Content
                value="request"
                class="min-h-0 flex-1 overflow-auto"
            >
                <MessageView message={flow.request} />
            </Tabs.Content>
            <Tabs.Content
                value="response"
                class="min-h-0 flex-1 overflow-auto"
            >
                {#if flow.response}
                    <MessageView message={flow.response} />
                {:else}
                    <p class="text-surface-600-400 text-sm">
                        {flow.summary.error || "No response"}
                    </p>
                {/if}
            </Tabs.Content>
        </Tabs>
    {/if}
</aside>
