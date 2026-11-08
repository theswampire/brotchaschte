<script lang="ts">
    import { createMutation, useQueryClient } from "@tanstack/svelte-query";
    import { isAxiosError } from "axios";
    import { api } from "@/lib/api/api";
    import type { GetOrdersRes } from "../api/orders";

    const queryClient = useQueryClient();
    const syncOrdersMutation = createMutation(() => ({
        mutationFn: api.orders.pullUpstreamOrders,
        onSuccess() {
            queryClient.invalidateQueries({ queryKey: ["orders"] });
            queryClient.invalidateQueries({ queryKey: ["sorted-orders"] });
        },
        onError(error) {
            if (isAxiosError<GetOrdersRes>(error)) {
                if (error.response?.data) {
                    console.log(error.response?.data);
                }
            }
        },
    }));
    async function syncOrders() {
        syncOrdersMutation.mutate();
    }
</script>

<button
    disabled={syncOrdersMutation.isPending}
    class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-300"
    on:click={() => syncOrders()}
    type="button"
>
    {#if syncOrdersMutation.isPending}
        Importieren...
    {/if}
    {#if !syncOrdersMutation.isPending}
        Importieren
    {/if}
</button>
