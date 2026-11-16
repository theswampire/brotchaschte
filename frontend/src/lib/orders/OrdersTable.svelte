<script lang="ts">
  import { createQuery } from "@tanstack/svelte-query";
  import { api } from "@/lib/api/api";
  import type { Order } from "@/lib/api/orders";
  import { navigate } from "sv-router/generated";

  interface Props {
    edit_order_callback: (order: Order) => void;
  }
  const { edit_order_callback }: Props = $props();

  const query = createQuery(() => ({
    queryKey: ["orders"],
    queryFn: api.orders.getOrders,
  }));
  // TODO: edit orders and manually geocode them
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
      <thead class="bg-neutral-300 text-sm text-slate-700">
        <tr class="whitespace-nowrap text-left">
          <th>Id</th>
          <th>Name</th>
          <th>E-Mail</th>
          <th>Tel</th>
          <th>Addresse</th>
          <th>PLZ</th>
          <th>Stadt</th>
          <th>#S</th>
          <th>#M</th>
          <th>#L</th>
          <th>Lieferzeit</th>
          <th>Bemerkungen</th>
        </tr>
      </thead>
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
            <td class="block max-w-20 overflow-hidden text-ellipsis px-2 py-1"
              >{o.id}</td
            >
            <td>{o.firstname} {o.lastname}</td>
            <td>{o.email}</td>
            <td>{o.phonenumber}</td>
            <td>{o.address}</td>
            <td>{o.postalcode}</td>
            <td>{o.city}</td>
            <td>{o.amountsizeS}</td>
            <td>{o.amountsizeM}</td>
            <td>{o.amountsizeL}</td>
            <td>{o.deliverymethod}</td>
            <td>{o.remarks}</td>
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
      <thead class="bg-neutral-300 text-sm text-slate-700">
        <tr class="whitespace-nowrap text-left">
          <th>Id</th>
          <th>Name</th>
          <th>E-Mail</th>
          <th>Tel</th>
          <th>Addresse</th>
          <th>PLZ</th>
          <th>Stadt</th>
          <th>#S</th>
          <th>#M</th>
          <th>#L</th>
          <th>Lieferzeit</th>
          <th>Bemerkungen</th>
        </tr>
      </thead>
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
            <td class="block max-w-20 overflow-hidden text-ellipsis px-2 py-1"
              >{o.id}</td
            >
            <td>{o.firstname} {o.lastname}</td>
            <td>{o.email}</td>
            <td>{o.phonenumber}</td>
            <td>{o.address}</td>
            <td>{o.postalcode}</td>
            <td>{o.city}</td>
            <td>{o.amountsizeS}</td>
            <td>{o.amountsizeM}</td>
            <td>{o.amountsizeL}</td>
            <td>{o.deliverymethod}</td>
            <td>{o.remarks}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <h2 class="block my-5 font-bold text-2xl">Bestellungen mit Lokalisierung</h2>
  <div class="w-full overflow-x-auto">
    <table class="w-full">
      <thead class="bg-neutral-300 text-sm text-slate-700">
        <tr class="whitespace-nowrap text-left">
          <th>Id</th>
          <th>Name</th>
          <th>E-Mail</th>
          <th>Tel</th>
          <th>Addresse</th>
          <th>PLZ</th>
          <th>Stadt</th>
          <th>#S</th>
          <th>#M</th>
          <th>#L</th>
          <th>Lieferzeit</th>
          <th>Bemerkungen</th>
        </tr>
      </thead>
      <tbody>
        {#each query.data.data.geocoded_orders as o}
          <tr
            class="whitespace-nowrap even:bg-blue-200 hover:bg-blue-50 hover:even:bg-blue-100 transition-colors"
            onclick={() => {
              edit_order_callback(o);
            }}
          >
            <td class="block max-w-20 overflow-hidden text-ellipsis px-2 py-1"
              >{o.id}</td
            >
            <td>{o.firstname} {o.lastname}</td>
            <td>{o.email}</td>
            <td>{o.phonenumber}</td>
            <td>{o.address}</td>
            <td>{o.postalcode}</td>
            <td>{o.city}</td>
            <td>{o.amountsizeS}</td>
            <td>{o.amountsizeM}</td>
            <td>{o.amountsizeL}</td>
            <td>{o.deliverymethod}</td>
            <td>{o.remarks}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}
