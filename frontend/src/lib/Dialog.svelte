<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children?: Snippet;
    show: boolean;
    onclose?: () => void;
    class?: string;
  }
  let {
    children,
    show = $bindable(),
    onclose,
    class: className,
  }: Props = $props();

  let dialogEl: HTMLDialogElement;

  $effect(() => {
    if (!dialogEl) return;
    if (show && !dialogEl.open) {
      dialogEl.showModal();
    } else if (!show && dialogEl.open) {
      dialogEl.close();
    }
  });

  function handleClose() {
    show = false;
    if (onclose) onclose();
  }
</script>

<dialog bind:this={dialogEl} onclose={handleClose} class={className}>
  {@render children?.()}
</dialog>
