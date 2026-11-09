<script lang="ts">
    import { createQuery } from "@tanstack/svelte-query";
    import { api } from "@/lib/api/api";
    import CreateDeliveryForm from "./CreateDeliveryForm.svelte";
    import { p } from "sv-router/generated";

    const query = createQuery(() => ({
        queryKey: ["deliveries"],
        queryFn: api.deliveries.getDeliveries,
    }));
</script>

<section
    class="bg-white left-1/2 top-0 w-1/2 h-full shadow border-l-gray-300 absolute p-4"
>
    <h2 class="font-bold text-3xl">Lieferungen</h2>

    <CreateDeliveryForm />

    {#if query.isPending}
        <p>Loading</p>
    {/if}
    {#if query.isSuccess}
        {#each query.data.data.deliveries as delivery}
            <a
                class="flex gap-3"
                href={p("/deliveries/:deliveryId", {
                    params: { deliveryId: delivery.id.toString() },
                })}
                data-view-transition
            >
                <div
                    class="aspect-square w-10"
                    style:background-color={delivery.color}
                ></div>
                {delivery.id}
                {delivery.label}
                {delivery.until}
            </a>
        {/each}
    {/if}
</section>
