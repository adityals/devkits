import {
  registerCompletions,
  type Completion,
  type CompletionSource,
} from "prism-code-editor/autocomplete";
import {
  completeKeywords,
  completeScope,
  jsContext,
  type JSContext,
} from "prism-code-editor/autocomplete/javascript";

// Mirrors flowObject and headerObject in internal/mitm/rule/hooks.go; change both together.
const members: Record<string, Completion[]> = {
  flow: [
    { label: "url", detail: "string", icon: "property" },
    { label: "host", detail: "string", icon: "property" },
    { label: "path", detail: "string", icon: "property" },
    { label: "query", detail: "Record<string, string[]>", icon: "property" },
  ],
  h: [
    { label: "get", detail: "(name: string) => string", icon: "function" },
    { label: "values", detail: "(name: string) => string[]", icon: "function" },
    { label: "set", detail: "(name: string, value: string) => void", icon: "function" },
    { label: "add", detail: "(name: string, value: string) => void", icon: "function" },
    { label: "del", detail: "(name: string) => void", icon: "function" },
    { label: "toObject", detail: "() => Record<string, string[]>", icon: "function" },
  ],
};

const hook = (name: string, args: string, returns: string): Completion => {
  const insert = `${name}(${args}) {},`;

  return {
    label: name,
    detail: `(${args}) => ${returns}`,
    icon: "snippet",
    insert,
    // Cursor lands between the braces.
    tabStops: [insert.length - 2],
  };
};

const ruleKeys: Completion[] = [
  { label: "match", detail: "string", icon: "property" },
  { label: "modifier", detail: '"=" | "^~" | "~" | "~*"', icon: "property" },
  { label: "host", detail: "string | string[]", icon: "property" },
  hook("onRequestHeader", "flow, h", "void"),
  hook("onResponseHeader", "flow, h", "void"),
  hook("onRequestBody", "flow, body", "string | undefined"),
  hook("onResponseBody", "flow, body", "string | undefined"),
];

const ruleSource: CompletionSource<JSContext> = ({ path, pos }) => {
  if (!path) {
    return null;
  }

  const options = path.length === 1 ? ruleKeys : members[path.slice(0, -1).join(".")];
  if (!options) {
    return null;
  }

  return {
    from: pos - path[path.length - 1].length,
    options,
  };
};

const modifiers: Completion[] = [
  { label: "=", detail: "exact path", icon: "constant" },
  { label: "^~", detail: "prefix; skips regex when longest", icon: "constant" },
  { label: "~", detail: "regex", icon: "constant" },
  { label: "~*", detail: "case-insensitive regex", icon: "constant" },
];

// jsContext disables completion inside strings, so this source reads the line instead of path.
const modifierSource: CompletionSource = ({ lineBefore, pos }) => {
  const typed = /\bmodifier\s*:\s*["'`]([^"'`]*)$/.exec(lineBefore)?.[1];
  if (typed === undefined) {
    return null;
  }

  return {
    from: pos - typed.length,
    options: modifiers,
  };
};

registerCompletions(["javascript"], {
  context: jsContext,
  sources: [
    ruleSource,
    modifierSource,
    completeScope({
      console,
      body: "",
    }),
    completeKeywords,
  ],
});
