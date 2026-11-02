<script lang="ts">
    import { onMount } from "svelte";
    import { Events } from "@wailsio/runtime";
    import { AppBar, Switch } from "@skeletonlabs/skeleton-svelte";
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

<AppBar class="sticky top-0 z-10 border-b border-surface-200-800">
    <AppBar.Toolbar class="grid-cols-[auto_1fr_auto] items-center">
        <AppBar.Lead>
            <h1 class="h6">HTTP Proxy</h1>
        </AppBar.Lead>
        <AppBar.Headline>
            {#if server.address}
                <code class="code text-xs">{server.address}</code>
            {/if}
        </AppBar.Headline>
        <AppBar.Trail class="items-center gap-3">
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
        </AppBar.Trail>
    </AppBar.Toolbar>
    {#if server.error}
        <p class="preset-tonal-error px-4 py-2 text-sm">
            {server.error}
        </p>
    {/if}
</AppBar>
