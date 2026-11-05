<script lang="ts">
    import { onMount } from "svelte";
    import { Switch } from "@skeletonlabs/skeleton-svelte";
    import { Plus, RefreshCw, Trash } from "@lucide/svelte";
    import { RuleList } from "./list.svelte";
    import { toaster } from "../../../lib/toaster";

    const list = new RuleList();
    let newId = $state("");
    let confirming = $state<string | null>(null);

    const create = async (event: SubmitEvent) => {
        event.preventDefault();

        const created = await list.create(newId.trim());
        if (created) {
            newId = "";
        }
    };

    const handleReload = async () => {
        const reloaded = await list.load();
        if (reloaded) {
            toaster.success({
                title: "Rules reloaded",
                description: `${list.rules.length} rules from disk`,
            });
        }
    };

    const handleDelete = (id: string) => {
        if (confirming !== id) {
            confirming = id;
            return;
        }

        confirming = null;
        list.remove(id);
    };

    onMount(() => {
        list.load();
    });
</script>

<div class="flex h-full flex-col gap-3 pt-3">
    <form class="flex gap-2" onsubmit={create}>
        <input class="input" placeholder="new-rule-id" bind:value={newId} />
        <button
            type="submit"
            class="btn preset-filled-primary-500"
            disabled={newId.trim() === ""}
        >
            <Plus size={16} />
            <span>New</span>
        </button>
        <button
            type="button"
            class="btn-icon preset-tonal"
            aria-label="Reload from disk"
            title="Reload from disk"
            onclick={handleReload}
        >
            <RefreshCw size={16} />
        </button>
        <span
            class="badge preset-tonal ml-auto self-center whitespace-nowrap tabular-nums"
        >
            {list.enabled} / {list.capacity} enabled
        </span>
    </form>
    {#if list.error}
        <p class="card preset-tonal-error p-3 text-sm">{list.error}</p>
    {/if}
    <ul
        class="card preset-outlined-surface-200-800 divide-surface-200-800 min-h-0 flex-1 divide-y overflow-auto"
    >
        {#each list.rules as rule (rule.id)}
            <li class="flex items-center gap-3 px-4 py-2">
                <code class="code flex-1 truncate">{rule.id}</code>
                <Switch
                    checked={rule.enabled}
                    onCheckedChange={() => list.toggle(rule)}
                >
                    <Switch.Control
                        class="data-[state=checked]:preset-filled-success-500"
                    >
                        <Switch.Thumb />
                    </Switch.Control>
                    <Switch.Label class="sr-only">Enable {rule.id}</Switch.Label
                    >
                    <Switch.HiddenInput />
                </Switch>
                <button
                    type="button"
                    class={[
                        "btn btn-sm",
                        confirming === rule.id
                            ? "preset-filled-error-500"
                            : "preset-tonal",
                    ]}
                    aria-label={`Delete ${rule.id}`}
                    onclick={() => handleDelete(rule.id)}
                    onblur={() => (confirming = null)}
                >
                    <Trash size={16} />
                    {#if confirming === rule.id}
                        <span>Confirm</span>
                    {/if}
                </button>
            </li>
        {:else}
            <li class="text-surface-600-400 p-4 text-sm">No rules yet.</li>
        {/each}
    </ul>
</div>
