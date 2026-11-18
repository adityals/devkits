import { convert, sources, targetsFor } from "./provider/registry";
import type { Lang } from "./provider/types";

export class Converter {
  input = $state("");
  output = $state("");
  error = $state("");
  from = $state<Lang>(sources[0]);
  to = $state<Lang>(targetsFor(sources[0])[0]);
  typeName = $state("Root");
  goPackage = $state("main");
  targets = $derived(targetsFor(this.from));
  #latest = 0;

  async run(): Promise<void> {
    if (!this.input.trim()) {
      this.output = "";
      this.error = "";
      return;
    }

    // Only the newest run may write, so a slow run can't overwrite a newer result.
    const run = ++this.#latest;
    try {
      const output = await convert(this.input, this.from, this.to, {
        typeName: this.typeName,
        goPackage: this.goPackage,
      });

      if (run === this.#latest) {
        this.output = output;
        this.error = "";
      }
    } catch (err) {
      if (run === this.#latest) {
        this.error = err instanceof Error ? err.message : String(err);
      }
    }
  }
}
