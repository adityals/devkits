import { Lang as BackendLang } from "../../../../bindings/devkits/internal/typegen/convert";
import { Convert } from "../../../../bindings/devkits/internal/typegen/convert/convertservice";
import type { Lang, Provider } from "./types";

const langs: Partial<Record<Lang, BackendLang>> = {
  go: BackendLang.LangGo,
  ts: BackendLang.LangTS,
};

export class BackendProvider implements Provider {
  async convert(
    input: string,
    from: Lang,
    to: Lang,
  ): Promise<string> {
    const source = langs[from];
    const target = langs[to];
    if (!source || !target) {
      throw new Error(`backend: unsupported ${from} → ${to}`);
    }

    return Convert(
      source,
      target,
      input,
    );
  }
}
