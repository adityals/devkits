<script lang="ts">
    import Router from "svelte-spa-router";
    import Shell from "./layout/Shell.svelte";
    import routes from "./routes";
</script>

<Shell>
    <Router {routes} />
</Shell>
