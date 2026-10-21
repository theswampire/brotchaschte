<script lang="ts">
    import { createMutation, useQueryClient } from "@tanstack/svelte-query";
    import { api } from "@/lib/api/api";
    import type { GetOrdersRes } from "../api/orders";
    import { isAxiosError } from "axios";

    const queryClient = useQueryClient();
    const generateOrdersMutation = createMutation(() => ({
        mutationFn: api.orders.generateRandomOrders,
        onSuccess() {
            queryClient.invalidateQueries({ queryKey: ["orders"] });
        },
        onError(error) {
            if (isAxiosError<GetOrdersRes>(error)) {
                if (error.response?.data) {
                    console.log(error["response"]?.data);
                }
            }
        },
    }));
    async function generate() {
        generateOrdersMutation.mutate();
    }
</script>

<button
    disabled={generateOrdersMutation.isPending}
    class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow"
    on:click={() => generate()}
    type="button">Zufällig Generieren</button
>
