<script lang="ts">
    import type { PrismEditor } from "prism-code-editor";
    import { autoComplete, fuzzyFilter } from "prism-code-editor/autocomplete";
    import autocompleteStyles from "prism-code-editor/autocomplete.css?inline";
    import autocompleteIcons from "prism-code-editor/autocomplete-icons.css?inline";
    import "prism-code-editor/prism/languages/javascript";
    import CodeEditor from "../../../lib/components/CodeEditor.svelte";
    import "./completions";

    type Props = {
        value: string;
        onChange: (value: string) => void;
        onSave: () => void;
    };

    const { value, onChange, onSave }: Props = $props();

    const addAutocomplete = (editor: PrismEditor, root: ShadowRoot) => {
        editor.addExtensions(autoComplete({ filter: fuzzyFilter }));

        // The editor renders in a shadow root, so page-level CSS never reaches the popup.
        const style = document.createElement("style");
        style.textContent = autocompleteStyles + autocompleteIcons;
        root.append(style);
    };

    const handleKeydown = (event: KeyboardEvent) => {
        if ((event.metaKey || event.ctrlKey) && event.key === "s") {
            event.preventDefault();
            onSave();
        }
    };
</script>

<div
    class="h-full min-h-0"
    onkeydown={handleKeydown}
    role="presentation"
>
    <CodeEditor
        {value}
        language="javascript"
        {onChange}
        setup={addAutocomplete}
    />
</div>
