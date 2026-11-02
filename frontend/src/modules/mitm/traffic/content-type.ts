import beautify from "js-beautify";
import type { Header } from "../../../../bindings/net/http";
import type { Lang } from "../../../lib/components/CodeBlock.svelte";

const LANGS: Record<string, Lang> = {
  "application/json": "json",
  "text/json": "json",
  "text/html": "html",
  "application/xhtml+xml": "html",
  "application/xml": "xml",
  "text/xml": "xml",
  "application/javascript": "javascript",
  "text/javascript": "javascript",
  "text/css": "css",
};

export const headerValue = (header: Header, name: string): string => header?.[name]?.[0] ?? "";

export const langFor = (contentType: string): Lang => {
  const mime = contentType.split(";")[0].trim().toLowerCase();

  if (LANGS[mime]) {
    return LANGS[mime];
  }

  if (mime.endsWith("+json")) {
    return "json";
  }

  if (mime.endsWith("+xml")) {
    return "xml";
  }

  return "text";
};

// Formatting runs on the UI thread, so large bodies are shown as received.
const MAX_FORMAT = 100_000;

const beautifiers: Partial<Record<Lang, (body: string) => string>> = {
  html: (body) => beautify.html(body, { indent_size: 2 }),
  xml: (body) => beautify.html(body, { indent_size: 2 }),
  css: (body) => beautify.css(body, { indent_size: 2 }),
  javascript: (body) => beautify.js(body, { indent_size: 2 }),
};

export const formatBody = (body: string, lang: Lang): string => {
  if (body.length > MAX_FORMAT) {
    return body;
  }

  if (lang === "json") {
    try {
      return JSON.stringify(JSON.parse(body), null, 2);
    } catch {
      return body;
    }
  }

  return beautifiers[lang]?.(body) ?? body;
};
