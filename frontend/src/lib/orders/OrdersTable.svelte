<script lang="ts">
  import { createQuery } from "@tanstack/svelte-query";
  import { api } from "@/lib/api/api";
  import type { Order } from "@/lib/api/orders";
  import { navigate } from "sv-router/generated";
  import OrderTableHeaders from "./OrderTableHeaders.svelte";
  import OrderTableRow from "./OrderTableRow.svelte";

  interface Props {
    edit_order_callback: (order: Order) => void;
  }
  const { edit_order_callback }: Props = $props();

  const query = createQuery(() => ({
    queryKey: ["orders"],
    queryFn: api.orders.getOrders,
  }));
</script>

{#if query.isPending}
  <p>Loading</p>
{/if}
{#if query.isSuccess}
  <h2 class="block my-5 font-bold text-2xl">
    Bestellungen mit fehlgeschlagener Lokalisierung
  </h2>
  <div class="w-full overflow-x-auto">
    <table class="w-full">
      <OrderTableHeaders />
      <tbody>
        {#each query.data.data.failed_geocoded_orders as o}
          <tr
            class="whitespace-nowrap even:bg-blue-200 hover:bg-blue-50 hover:even:bg-blue-100 transition-colors"
            onclick={() => {
              navigate("/orders/:orderId/manual", {
                params: {
                  orderId: o.id,
                },
                viewTransition: true,
              });
            }}
          >
            <OrderTableRow order={o} />
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <h2 class="block my-5 font-bold text-2xl">
    Bestellungen (nicht lokalisiert)
  </h2>
  <div class="w-full overflow-x-auto">
    <table class="w-full">
      <OrderTableHeaders />
      <tbody>
        {#each query.data.data.not_geocoded_orders as o}
          <tr
            class="whitespace-nowrap even:bg-blue-200 hover:bg-blue-50 hover:even:bg-blue-100 transition-colors"
            onclick={() => {
              navigate("/orders/:orderId/manual", {
                params: {
                  orderId: o.id,
                },
                viewTransition: true,
              });
            }}
          >
            <OrderTableRow order={o} />
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <h2 class="block my-5 font-bold text-2xl">Bestellungen mit Lokalisierung</h2>
  <div class="w-full overflow-x-auto">
    <table class="w-full">
      <OrderTableHeaders />
      <tbody>
        {#each query.data.data.geocoded_orders as o}
          <tr
            class="whitespace-nowrap even:bg-blue-200 hover:bg-blue-50 hover:even:bg-blue-100 transition-colors"
            onclick={() => {
              edit_order_callback(o);
            }}
          >
            <OrderTableRow order={o} />
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}
