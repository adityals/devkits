import wrap from "svelte-spa-router/wrap";
import MitmPage from "./routes/MitmPage.svelte";

const routes = {
  "/mitm": MitmPage,
  "*": wrap({ asyncComponent: () => import("./routes/NotFound.svelte") }),
};

export default routes;
