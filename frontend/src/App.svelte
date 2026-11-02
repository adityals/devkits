<script lang="ts">
    import Router, { replace, router } from "svelte-spa-router";
    import Shell from "./layout/Shell.svelte";
    import routes from "./routes";

    $effect(() => {
        if (router.location === "/") {
            replace("/mitm");
        }
    });
</script>

<Shell>
    <Router {routes} />
</Shell>
