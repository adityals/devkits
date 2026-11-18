<script lang="ts">
    import { untrack } from "svelte";
    import type { PrismEditor } from "prism-code-editor";
    import { basicEditor } from "prism-code-editor/setups";

    type Props = {
        value: string;
        // The grammar must be imported by the caller, e.g. "prism-code-editor/prism/languages/json".
        language: string;
        onChange: (value: string) => void;
        // Runs once after creation; the shadow root is where extension styles must go.
        setup?: (editor: PrismEditor, root: ShadowRoot) => void;
    };

    const {
        value,
        language,
        onChange,
        setup,
    }: Props = $props();
    let editor = $state.raw<PrismEditor>();

    // shortcut: theme is picked once from the OS setting, add a listener if live switching matters.
    const theme = window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "github-dark"
        : "github-light";

    const editorComponent = (node: HTMLElement) => {
        // untracked: reading props here would re-run the attachment and rebuild the editor on every keystroke.
        const created = untrack(() =>
            basicEditor(node, {
                language,
                theme,
                value,
                onUpdate: (next) => onChange(next),
            }),
        );
        editor = created;

        const root = node.shadowRoot;
        if (root) {
            untrack(() => setup?.(created, root));
        }

        return () => created.remove();
    };

    // Only push outside changes, like a reload from disk; pushing on every keystroke would reset the cursor.
    $effect(() => {
        if (editor && value !== editor.value) {
            editor.setOptions({ value });
        }
    });
</script>

<div
    class="grid h-full min-h-0 overflow-auto"
    {@attach editorComponent}
></div>
