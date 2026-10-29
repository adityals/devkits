<script lang="ts">
    import { onMount } from "svelte";
    import { Events } from "@wailsio/runtime";
    import { Switch } from "@skeletonlabs/skeleton-svelte";
    import {
        Start,
        Stop,
    } from "../../../bindings/devkits/internal/proxy/proxyservice";
    import { ProxyStatus } from "../../../bindings/events";

    let server = $state({ running: false, error: "", address: "" });

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

        const unsubProxyStatusEvt = Events.On(ProxyStatus, (status) => {
            server.running = status.data.running;
            server.error = status.data.error ?? "";
            server.address = status.data.address;
        });

        return () => {
            unsubProxyStatusEvt();
        };
    });
</script>

<div class="flex min-h-screen flex-col">
    <main class="mx-auto w-full max-w-2xl space-y-4 p-6">
        <section class="card preset-outlined-surface-200-800 space-y-4 p-5">
            <header class="flex items-start justify-between gap-4">
                <div class="space-y-1">
                    <h2 class="h5">HTTP Proxy</h2>
                    <p class="text-surface-600-400 text-sm">
                        {#if server.address}
                            Intercepts traffic on <code class="code"
                                >{server.address}</code
                            >
                        {/if}
                    </p>
                </div>
                <Switch
                    checked={server.running}
                    onCheckedChange={(detail) => {
                        detail.checked
                            ? void handleStartProxyServer()
                            : void handleStopProxyServer();
                    }}
                >
                    <Switch.Control
                        class="data-[state=checked]:preset-filled-success-500"
                    >
                        <Switch.Thumb />
                    </Switch.Control>
                    <Switch.Label class="sr-only">Toggle proxy</Switch.Label>
                    <Switch.HiddenInput />
                </Switch>
            </header>
            <hr class="hr" />
            <dl class="grid grid-cols-2 gap-4 text-sm">
                <div>
                    <dt class="text-surface-600-400">Status</dt>
                    <dd class="font-semibold">
                        {server.running ? "Running" : "Stopped"}
                    </dd>
                </div>
                <div>
                    <dt class="text-surface-600-400">Address</dt>
                    <dd class="font-semibold">{server.address}</dd>
                </div>
            </dl>
            {#if server.error}
                <p class="card preset-tonal-error p-3 text-sm">
                    {server.error}
                </p>
            {/if}
        </section>
    </main>
</div>
