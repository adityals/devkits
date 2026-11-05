import {
  Capacity,
  Create,
  Delete,
  List,
  Reload,
  Update,
} from "../../../../bindings/devkits/internal/mitm/rule/ruleservice";
import type { Rule } from "../../../../bindings/devkits/internal/mitm/rule";

const TEMPLATE = `const rule = {
  match: "/",
  onRequestHeader(h) {},
  onResponseHeader(h) {},
};
`;

export class RuleList {
  rules = $state.raw<Rule[]>([]);
  error = $state("");
  capacity = $state(0);
  enabled = $derived(this.rules.filter((r) => r.enabled).length);

  // load goes through Reload so edits made outside the app show up.
  load(): Promise<boolean> {
    return this.#run(async () => {
      this.capacity = await Capacity();
      await Reload();
    });
  }

  toggle(rule: Rule): Promise<boolean> {
    return this.#run(() => Update({ ...rule, enabled: !rule.enabled }));
  }

  // New rules start disabled, so they never hit the limit or touch traffic before they are written.
  create(id: string): Promise<boolean> {
    return this.#run(() => Create({ id, enabled: false, script: TEMPLATE }));
  }

  remove(id: string): Promise<boolean> {
    return this.#run(() => Delete(id));
  }

  // #run always re-lists: some failures, like going over the limit, still change what is on disk.
  async #run(op: () => Promise<unknown>): Promise<boolean> {
    this.error = "";
    let ok = true;

    try {
      await op();
    } catch (error) {
      this.error = String(error);
      ok = false;
    }

    try {
      this.rules = (await List()) ?? [];
    } catch (error) {
      this.error = String(error);
      ok = false;
    }

    return ok;
  }
}
