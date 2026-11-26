<script lang="ts">
    import { onMount } from "svelte";
    import { Switch } from "@skeletonlabs/skeleton-svelte";
    import { Pencil, Plus, RefreshCw, Save, Trash } from "@lucide/svelte";
    import type { Rule } from "../../../../bindings/devkits/internal/mitm/rule";
    import { RuleList } from "./list.svelte";
    import { toaster } from "../../../lib/toaster";
    import Loading from "../../../lib/components/Loading/Loading.svelte";

    const list = new RuleList();

    let newId = $state("");
    let confirming = $state<string | null>(null);
    let selectedId = $state<string | null>(null);
    let draft = $state("");
    let renameId = $state("");
    let discarding = $state<string | null>(null);
    let editorModule = $state<Promise<typeof import("./RuleEditor.svelte")>>();

    const selected = $derived(list.rules.find((r) => r.id === selectedId));
    const dirty = $derived(selected !== undefined && draft !== selected.script);

    const handleSelect = (rule: Rule) => {
        if (rule.id === selectedId) {
            return;
        }

        if (dirty && discarding !== rule.id) {
            discarding = rule.id;
            toaster.warning({
                title: "Unsaved changes",
                description: `Save ${selectedId}, or click ${rule.id} again to discard.`,
            });
            return;
        }

        discarding = null;
        selectedId = rule.id;
        draft = rule.script;
        renameId = rule.id;
        editorModule ??= import("./RuleEditor.svelte");
    };

    // Unsaved edits survive a rename: the draft is kept and still compared against the renamed rule's script.
    const handleRename = async (event: SubmitEvent) => {
        event.preventDefault();

        if (!selected) {
            return;
        }

        const newId = renameId.trim();
        const renamed = await list.rename(selected.id, newId);
        if (renamed) {
            selectedId = newId;
        }
    };

    const handleSave = async () => {
        if (!selected || !dirty) {
            return;
        }

        const saved = await list.save(selected, draft);
        if (saved) {
            discarding = null;
            toaster.success({ title: `${selected.id} saved` });
        }
    };

    const create = async (event: SubmitEvent) => {
        event.preventDefault();

        const created = await list.create(newId.trim());
        if (created) {
            newId = "";
        }
    };

    // shortcut: a reload keeps unsaved edits and the next save overwrites the disk version; add a conflict prompt if that bites.
    const handleReload = async () => {
        const wasClean = !dirty;
        const reloaded = await list.load();
        if (wasClean) {
            draft = selected?.script ?? "";
        }

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
        <input
            class="input"
            type="text"
            placeholder="new-rule-id"
            autocapitalize="none"
            bind:value={newId}
        />
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
    <div class="grid min-h-0 flex-1 grid-cols-[18rem_1fr] gap-3">
        <ul
            class="card preset-outlined-surface-200-800 divide-surface-200-800 min-h-0 divide-y overflow-auto"
        >
            {#each list.rules as rule (rule.id)}
                <li
                    class={[
                        "flex items-center gap-3 px-4 py-2",
                        selectedId === rule.id && "preset-tonal-primary",
                    ]}
                >
                    <button
                        type="button"
                        class="flex-1 truncate text-left"
                        onclick={() => handleSelect(rule)}
                    >
                        <code class="code">{rule.id}</code>
                    </button>
                    <Switch
                        checked={rule.enabled}
                        onCheckedChange={() => list.toggle(rule)}
                    >
                        <Switch.Control
                            class="data-[state=checked]:preset-filled-success-500"
                        >
                            <Switch.Thumb />
                        </Switch.Control>
                        <Switch.Label class="sr-only"
                            >Enable {rule.id}</Switch.Label
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
        <section
            class="card preset-outlined-surface-200-800 flex min-h-0 flex-col overflow-hidden"
        >
            {#if selected}
                <header
                    class="border-surface-200-800 flex items-center gap-3 border-b px-4 py-2"
                >
                    <form
                        class="flex items-center gap-2"
                        onsubmit={handleRename}
                    >
                        <input
                            type="text"
                            autocapitalize="none"
                            class="input w-48 font-mono text-sm"
                            aria-label="Rule ID"
                            bind:value={renameId}
                        />
                        <span class="text-surface-600-400 text-sm">.js</span>
                        {#if renameId.trim() !== selected.id}
                            <button
                                type="submit"
                                class="btn btn-sm preset-tonal"
                                disabled={renameId.trim() === ""}
                            >
                                <Pencil size={16} />
                                <span>Rename</span>
                            </button>
                        {/if}
                    </form>
                    {#if dirty}
                        <span class="badge preset-tonal-warning">Unsaved</span>
                    {/if}
                    <button
                        type="button"
                        class="btn btn-sm preset-filled-primary-500 ml-auto"
                        disabled={!dirty}
                        onclick={handleSave}
                    >
                        <Save size={16} />
                        <span>Save</span>
                    </button>
                </header>
                <div class="min-h-0 flex-1">
                    {#await editorModule}
                        <Loading class="p-3" />
                    {:then module}
                        {#if module}
                            <module.default
                                value={draft}
                                onChange={(next) => (draft = next)}
                                onSave={handleSave}
                            />
                        {/if}
                    {:catch error}
                        <p class="card preset-tonal-error m-3 p-3 text-sm">
                            {String(error)}
                        </p>
                    {/await}
                </div>
            {:else}
                <p class="text-surface-600-400 p-4 text-sm">
                    Select a rule to edit it.
                </p>
            {/if}
        </section>
    </div>
</div>
