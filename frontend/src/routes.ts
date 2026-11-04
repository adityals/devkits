import wrap from "svelte-spa-router/wrap";
import HomePage from "./routes/HomePage.svelte";

const routes = {
  "/": HomePage,
  "/mitm": wrap({ asyncComponent: () => import("./routes/MitmPage.svelte") }),
  "*": wrap({ asyncComponent: () => import("./routes/NotFound.svelte") }),
};

export default routes;
