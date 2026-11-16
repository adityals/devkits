<script lang="ts">
    import { untrack } from "svelte";
    import type { PrismEditor } from "prism-code-editor";
    import { basicEditor } from "prism-code-editor/setups";
    import "prism-code-editor/prism/languages/javascript";

    type Props = {
        value: string;
        onChange: (value: string) => void;
        onSave: () => void;
    };

    const { value, onChange, onSave }: Props = $props();
    let editor = $state.raw<PrismEditor>();

    // shortcut: theme is picked once from the OS setting, add a listener if live switching matters.
    const theme = window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "github-dark"
        : "github-light";

    const editorComponent = (node: HTMLElement) => {
        editor = basicEditor(node, {
            language: "javascript",
            theme,
            // untracked: reading value here would re-run the attachment and rebuild the editor on every keystroke.
            value: untrack(() => value),
            onUpdate: (next) => onChange(next),
        });

        return () => editor?.remove();
    };

    // Only push outside changes, like a reload from disk; pushing on every keystroke would reset the cursor.
    $effect(() => {
        if (editor && value !== editor.value) {
            editor.setOptions({ value });
        }
    });

    const handleKeydown = (event: KeyboardEvent) => {
        if ((event.metaKey || event.ctrlKey) && event.key === "s") {
            event.preventDefault();
            onSave();
        }
    };
</script>

<div
    class="grid h-full min-h-0 overflow-auto"
    onkeydown={handleKeydown}
    role="presentation"
    {@attach editorComponent}
></div>
