import { QuicktypeProvider } from "./quicktype";
import type { ConvertOptions, Lang, Provider, Route } from "./types";

const quicktype = new QuicktypeProvider();

export const routes: Partial<Record<Route, Provider>> = {
  json_go: quicktype,
  json_ts: quicktype,
};

const pairs = (Object.keys(routes) as Route[]).map(
  (route) => route.split("_") as [Lang, Lang],
);

export const sources: Lang[] = [...new Set(pairs.map(([from]) => from))];

export const targetsFor = (from: Lang): Lang[] =>
  pairs.filter(([source]) => source === from).map(([, to]) => to);

export const convert = (
  input: string,
  from: Lang,
  to: Lang,
  opts: ConvertOptions,
): Promise<string> => {
  const provider = routes[`${from}_${to}`];
  if (!provider) {
    return Promise.reject(new Error(`no provider for ${from} → ${to}`));
  }

  return provider.convert(input, from, to, opts);
};
