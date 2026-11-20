import type { LanguageName } from "quicktype-core";
import type { ConvertOptions, Lang, Provider } from "./types";

const targets: Partial<Record<Lang, LanguageName>> = {
  go: "go",
  ts: "typescript",
};

export class QuicktypeProvider implements Provider {
  async convert(input: string, from: Lang, to: Lang, opts: ConvertOptions): Promise<string> {
    const lang = targets[to];
    if (from !== "json" || !lang) {
      throw new Error(`quicktype: unsupported ${from} → ${to}`);
    }

    // Lazy: quicktype-core is ~4 MB and only this page needs it.
    const { InputData, jsonInputForTargetLanguage, quicktype } = await import("quicktype-core");

    const jsonInput = jsonInputForTargetLanguage(lang);
    await jsonInput.addSource({
      name: opts.typeName,
      samples: [input],
    });

    const inputData = new InputData();
    inputData.addInput(jsonInput);

    const result = await quicktype({
      inputData,
      lang,
      rendererOptions:
        to === "go"
          ? {
              "just-types-and-package": true,
              package: opts.goPackage,
            }
          : { "just-types": true },
    });

    return result.lines.join("\n");
  }
}
