<script lang="ts">
    import { onMount } from "svelte";
    import { Events } from "@wailsio/runtime";
    import {
        Start,
        Stop,
    } from "../../bindings/devkits/internal/proxy/proxyservice.js";

    const startProxyServer = async () => {
        try {
            await Start();
        } catch (error) {
            console.error("proxy start error:", error);
        }
    };

    const handleStopServer = async () => {
        try {
            await Stop();
        } catch (error) {
            console.error("proxy stop error:", error);
        }
    };

    onMount(() => {
        startProxyServer();

        const unsubProxyStatusEvt = Events.On("proxy:status", (status) => {
            console.log("proxy status:", status.data);
        });

        return () => {
            unsubProxyStatusEvt();
        };
    });
</script>

<h1>Home</h1>
<button on:click={handleStopServer}>Stop Server</button>
