<script lang="ts">
    import { createMutation, useQueryClient } from "@tanstack/svelte-query";
    import { isAxiosError } from "axios";
    import { api } from "@/lib/api/api";
    import type { GeocodeOrdersRes } from "../api/orders";

    const queryClient = useQueryClient();
    const geocodeOrdersMutation = createMutation(() => ({
        mutationFn: api.orders.geocodeOrders,
        onSuccess() {
            queryClient.invalidateQueries({ queryKey: ["orders"] });
            queryClient.invalidateQueries({ queryKey: ["sorted-orders"] });
        },
        onError(error) {
            if (isAxiosError<GeocodeOrdersRes>(error)) {
                if (error.response?.data) {
                    console.log(error.response?.data);
                }
            }
        },
    }));
    async function geocode() {
        geocodeOrdersMutation.mutate();
    }
</script>

<button
    disabled={geocodeOrdersMutation.isPending}
    class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-300"
    on:click={() => geocode()}
    type="button"
>
    {#if !geocodeOrdersMutation.isPending}
        Addressen lokalisieren
    {/if}
    {#if geocodeOrdersMutation.isPending}
        wird lokalisiert...
    {/if}
</button>
