import wrap from "svelte-spa-router/wrap";
import Home from "./routes/Home.svelte";

const routes = {
  "/": Home,
  "*": wrap({ asyncComponent: () => import("./routes/NotFound.svelte") }),
};

export default routes;
