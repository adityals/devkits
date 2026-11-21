import { convert, sources, targetsFor } from "./provider/registry";
import type { Lang } from "./provider/types";

export class Converter {
  input = $state("");
  output = $state("");
  error = $state("");
  loading = $state(false);
  from = $state<Lang>(sources[0]);
  to = $state<Lang>(targetsFor(sources[0])[0]);
  typeName = $state("Root");
  goPackage = $state("main");

  targets = $derived(targetsFor(this.from));

  #latest = 0;

  async run(): Promise<void> {
    const run = ++this.#latest;

    if (!this.input.trim()) {
      this.output = "";
      this.error = "";
      this.loading = false;
      return;
    }

    this.loading = true;
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
    } finally {
      if (run === this.#latest) {
        this.loading = false;
      }
    }
  }
}
