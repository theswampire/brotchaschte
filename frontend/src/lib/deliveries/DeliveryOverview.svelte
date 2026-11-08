<script lang="ts">
    interface Props {
        visible: boolean;
    }
    const { visible }: Props = $props();
</script>

<section
    class="bg-white left-1/2 top-0 w-1/2 h-full shadow border-l-gray-300 data-[visible=false]:hidden absolute p-4"
    data-visible={visible}
>
    <h2 class="font-bold text-3xl">Lieferungen</h2>
</section>
