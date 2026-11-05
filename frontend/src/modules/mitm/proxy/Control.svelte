<script lang="ts">
    import { onMount } from "svelte";
    import { Events } from "@wailsio/runtime";
    import { Switch } from "@skeletonlabs/skeleton-svelte";
    import { Globe, ShieldCheck, ShieldOff, Unplug } from "@lucide/svelte";
    import {
        Start,
        Stop,
    } from "../../../../bindings/devkits/internal/mitm/proxy/proxyservice";
    import {
        Remove,
        Status,
        Use,
    } from "../../../../bindings/devkits/internal/mitm/cert/certservice";
    import * as SysProxy from "../../../../bindings/devkits/internal/mitm/sysproxy/sysproxyservice";
    import { MitmProxyStatus } from "../../../../bindings/events";

    let server = $state({ running: false, error: "", address: "" });
    let ca = $state({ trusted: false, path: "", inProgress: false, error: "" });
    let sys = $state({ enabled: false, address: "", inProgress: false, error: "" });

    const refreshSys = async () => {
        const status = await SysProxy.Status();
        sys.enabled = status.enabled;
        sys.address = status.address;
    };

    const toggleSys = async () => {
        sys.inProgress = true;
        sys.error = "";

        try {
            await (sys.enabled ? SysProxy.Disable() : SysProxy.Enable());
            await refreshSys();
        } catch (error) {
            sys.error = String(error);
        } finally {
            sys.inProgress = false;
        }
    };

    const refreshCA = async () => {
        const status = await Status();
        ca.trusted = status.trusted;
        ca.path = status.path;
    };

    const toggleCA = async () => {
        ca.inProgress = true;
        ca.error = "";

        try {
            await (ca.trusted ? Remove() : Use());
            await refreshCA();
        } catch (error) {
            ca.error = String(error);
        } finally {
            ca.inProgress = false;
        }
    };

    const handleStartProxyServer = () => {
        Start().catch((error) => {
            server.error = String(error);
        });
    };

    const handleStopProxyServer = () => {
        Stop().catch((error) => {
            server.error = String(error);
        });
    };

    onMount(() => {
        handleStartProxyServer();
        refreshCA().catch((error) => {
            ca.error = String(error);
        });
        refreshSys().catch((error) => {
            sys.error = String(error);
        });

        const unsubProxyStatusEvt = Events.On(MitmProxyStatus, (status) => {
            server.running = status.data.running;
            server.error = status.data.error ?? "";
            server.address = status.data.address;

            // The backend turns the system proxy off when the proxy stops.
            if (!status.data.running) {
                refreshSys().catch((error) => {
                    sys.error = String(error);
                });
            }
        });

        return () => {
            unsubProxyStatusEvt();
        };
    });
</script>

<section class="card preset-outlined-surface-200-800 mx-4 mt-4 space-y-3 p-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center gap-3">
            <h1 class="h6">HTTP Proxy</h1>
            {#if server.address}
                <code class="code text-xs">{server.address}</code>
            {/if}
        </div>
        <div class="flex items-center gap-3">
            <span
                class={[
                    "badge text-xs",
                    server.running
                        ? "preset-tonal-success"
                        : "preset-tonal-surface",
                ]}
            >
                {server.running ? "Running" : "Stopped"}
            </span>
            <Switch
                checked={server.running}
                onCheckedChange={(detail) =>
                    detail.checked
                        ? void handleStartProxyServer()
                        : void handleStopProxyServer()}
            >
                <Switch.Control
                    class="data-[state=checked]:preset-filled-success-500"
                >
                    <Switch.Thumb />
                </Switch.Control>
                <Switch.Label class="sr-only">Toggle proxy</Switch.Label>
                <Switch.HiddenInput />
            </Switch>
        </div>
    </div>
    {#if server.error}
        <p class="card preset-tonal-error p-2 text-sm">
            {server.error}
        </p>
    {/if}
    <div
        class="border-surface-200-800 flex flex-wrap items-center justify-between gap-3 border-t pt-3"
    >
        <div class="flex min-w-0 items-center gap-3">
            <h2 class="text-sm font-semibold">HTTPS CA</h2>
            <span
                class={[
                    "badge text-xs",
                    ca.trusted
                        ? "preset-tonal-success"
                        : "preset-tonal-warning",
                ]}
            >
                {ca.trusted ? "Trusted" : "Untrusted"}
            </span>
            {#if ca.path}
                <code class="code truncate text-xs">{ca.path}</code>
            {/if}
        </div>
        <button
            type="button"
            class={[
                "btn btn-sm",
                ca.trusted ? "preset-tonal-error" : "preset-filled-primary-500",
            ]}
            disabled={ca.inProgress}
            onclick={toggleCA}
        >
            {#if ca.trusted}
                <ShieldOff size={16} />
                <span>Remove CA</span>
            {:else}
                <ShieldCheck size={16} />
                <span>Trust CA</span>
            {/if}
        </button>
    </div>
    {#if ca.error}
        <p class="card preset-tonal-error p-2 text-sm">
            {ca.error}
        </p>
    {/if}
    <div
        class="border-surface-200-800 flex flex-wrap items-center justify-between gap-3 border-t pt-3"
    >
        <div class="flex min-w-0 items-center gap-3">
            <h2 class="text-sm font-semibold">System Proxy</h2>
            <span
                class={[
                    "badge text-xs",
                    sys.enabled
                        ? "preset-tonal-success"
                        : "preset-tonal-surface",
                ]}
            >
                {sys.enabled ? "On" : "Off"}
            </span>
            {#if sys.address}
                <code class="code truncate text-xs">{sys.address}</code>
            {/if}
        </div>
        <button
            type="button"
            class={[
                "btn btn-sm",
                sys.enabled ? "preset-tonal-error" : "preset-filled-primary-500",
            ]}
            title={server.running ? undefined : "Start the proxy first"}
            disabled={sys.inProgress || (!sys.enabled && !server.running)}
            onclick={toggleSys}
        >
            {#if sys.enabled}
                <Unplug size={16} />
                <span>Disable</span>
            {:else}
                <Globe size={16} />
                <span>Enable</span>
            {/if}
        </button>
    </div>
    {#if sys.error}
        <p class="card preset-tonal-error p-2 text-sm">
            {sys.error}
        </p>
    {/if}
</section>
