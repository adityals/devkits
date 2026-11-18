import wrap from "svelte-spa-router/wrap";
import HomePage from "./routes/HomePage.svelte";
import RouteLoading from "./routes/RouteLoading.svelte";

const routes = {
  "/": HomePage,
  "/mitm": wrap({
    asyncComponent: () => import("./routes/MitmPage.svelte"),
    loadingComponent: RouteLoading,
  }),
  "/typegen": wrap({
    asyncComponent: () => import("./routes/TypegenPage.svelte"),
    loadingComponent: RouteLoading,
  }),
  "*": wrap({
    asyncComponent: () => import("./routes/NotFound.svelte"),
    loadingComponent: RouteLoading,
  }),
};

export default routes;
