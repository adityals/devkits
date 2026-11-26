<script lang="ts">
    import { onMount } from "svelte";
    import { Ban, Circle, Square } from "@lucide/svelte";
    import { Switch } from "@skeletonlabs/skeleton-svelte";
    import Loading from "../../../lib/components/Loading/Loading.svelte";
    import { TrafficLog } from "./log.svelte";

    const log = new TrafficLog();
    const recordLabel = $derived(log.recording ? "Stop recording" : "Record");
    let host = $state("");
    let path = $state("");
    let matchedOnly = $state(false);
    let selected = $state<string | null>(null);
    let detailModule = $state<Promise<typeof import("./Detail.svelte")>>();

    const rows = $derived.by(() => {
        const h = host.trim().toLowerCase();
        const p = path.trim().toLowerCase();

        return log.rows
            .filter(
                (r) =>
                    r.host.toLowerCase().includes(h) &&
                    r.path.toLowerCase().includes(p) &&
                    (!matchedOnly || Boolean(r.rule)),
            )
            .toReversed();
    });

    const statusClass = (status: number): string => {
        if (status >= 500 || status === 0) {
            return "text-error-500";
        }

        if (status >= 400) {
            return "text-warning-500";
        }

        return "text-success-500";
    };

    const handleSelect = (id: string) => {
        selected = id;
        detailModule ??= import("./Detail.svelte");
    };

    const handleRowKeydown = (event: KeyboardEvent, id: string) => {
        if (event.key === "Enter") {
            handleSelect(id);
        }
    };

    const formatSize = (bytes: number): string =>
        bytes < 1024 ? `${bytes} B` : `${(bytes / 1024).toFixed(1)} KB`;

    onMount(() => log.start());
</script>

<div class="flex h-full flex-col gap-3 pt-3">
    <div class="flex gap-2">
        <button
            type="button"
            class={[
                "btn-icon",
                log.recording ? "preset-filled-error-500" : "preset-tonal",
            ]}
            aria-label={recordLabel}
            title={recordLabel}
            onclick={() => log.toggleRecording()}
        >
            {#if log.recording}
                <Square size={16} />
            {:else}
                <Circle size={16} />
            {/if}
        </button>
        <button
            type="button"
            class="btn-icon preset-tonal"
            aria-label="Clear"
            title="Clear"
            disabled={log.rows.length === 0}
            onclick={() => {
                log.clear();
                selected = null;
            }}
        >
            <Ban size={16} />
        </button>
        <input
            class="input"
            type="search"
            placeholder="Filter host"
            bind:value={host}
        />
        <input
            class="input"
            type="search"
            placeholder="Filter path"
            bind:value={path}
        />
        <Switch
            class="shrink-0"
            checked={matchedOnly}
            onCheckedChange={(details) => (matchedOnly = details.checked)}
        >
            <Switch.Control
                class="data-[state=checked]:preset-filled-primary-500"
            >
                <Switch.Thumb />
            </Switch.Control>
            <Switch.Label class="text-sm whitespace-nowrap"
                >Matched</Switch.Label
            >
            <Switch.HiddenInput />
        </Switch>
        <span
            class="badge preset-tonal ml-auto self-center whitespace-nowrap tabular-nums"
            title="Records kept / capacity"
        >
            {log.rows.length} / {log.capacity}
        </span>
    </div>
    {#if log.error}
        <p class="card preset-tonal-error p-3 text-sm">{log.error}</p>
    {/if}
    <div class={["grid min-h-0 flex-1 gap-3", selected && "grid-cols-2"]}>
        <div class="table-wrap min-h-0 overflow-auto">
            <table class="table text-sm">
                <thead class="bg-surface-50-950 sticky top-0">
                    <tr>
                        <th>Time</th>
                        <th>Method</th>
                        <th>Status</th>
                        <th>Host</th>
                        <th>Path</th>
                        <th>Rule</th>
                        <th class="text-right">Duration</th>
                        <th class="text-right">Size</th>
                    </tr>
                </thead>
                <tbody>
                    {#each rows as row (row.id)}
                        <tr
                            title={row.error}
                            tabindex="0"
                            aria-selected={row.id === selected}
                            class={[
                                "cursor-pointer",
                                row.id === selected && "preset-tonal-primary",
                            ]}
                            onclick={() => handleSelect(row.id)}
                            onkeydown={(e) => handleRowKeydown(e, row.id)}
                        >
                            <td
                                >{new Date(
                                    row.startTime,
                                ).toLocaleTimeString()}</td
                            >
                            <td class="font-semibold">{row.method}</td>
                            <td class={statusClass(row.status)}>
                                {row.status || "ERR"}
                            </td>
                            <td>{row.host}</td>
                            <td class="max-w-md truncate">{row.path}</td>
                            <td>
                                {#if row.rule}
                                    <code class="code">{row.rule}</code>
                                {/if}
                            </td>
                            <td class="text-right">{row.durationMs} ms</td>
                            <td class="text-right">{formatSize(row.size)}</td>
                        </tr>
                    {:else}
                        <tr>
                            <td
                                colspan="8"
                                class="text-surface-600-400 text-center"
                            >
                                No traffic recorded
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
        {#if selected}
            {#await detailModule}
                <Loading class="p-3" />
            {:then module}
                {#if module}
                    <module.default
                        id={selected}
                        onclose={() => (selected = null)}
                    />
                {/if}
            {:catch error}
                <p class="card preset-tonal-error p-3 text-sm">
                    {String(error)}
                </p>
            {/await}
        {/if}
    </div>
</div>
