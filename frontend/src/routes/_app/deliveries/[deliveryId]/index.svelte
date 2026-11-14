<script lang="ts">
  import { api } from "@/lib/api/api";
  import type { Order } from "@/lib/api/orders";
  import EditDeliveryForm from "@/lib/deliveries/EditDeliveryForm.svelte";
  import Clock from "@/lib/icons/movingicons/clock.svelte";
  import Pencil from "@/lib/icons/movingicons/pencil.svelte";
  import X from "@/lib/icons/movingicons/x.svelte";
  import { getMarkerStoreContext } from "@/lib/map/markerstore.svelte";
  import { createMutation, createQuery } from "@tanstack/svelte-query";
  import { navigate, route } from "sv-router/generated";

  const params = route.getParams("/deliveries/:deliveryId");
  const deliveryId = Number.parseInt(params.deliveryId);

  const query = createQuery(() => ({
    queryKey: ["delivery", params.deliveryId],
    queryFn: () => api.deliveries.getDelivery(params.deliveryId),
  }));

  const assignOrdersMutation = createMutation(() => ({
    mutationFn: api.deliveries.assignOrdersToDelivery,
    onSuccess() {
      changed = false;
      markers.commitChanges();
    },
  }));

  let delivery = $derived(query.isSuccess ? query.data.data.delivery : null);
  // orders is not derived because we need to mutate state when adding / removing orders
  // which is simpler to do when we can directly mutate the array
  let orders = $state<Order[]>([]);

  const deselected_orders = new Set<string>();
  let changed = $state(false);
  let open_edit_dialog = $state(false);

  const markers = getMarkerStoreContext();
  markers.setToDeliveryEditMode(deliveryId);
  markers.setMarkerOnClick((id) => {
    const assignment = markers.getOrderAssignment(id);
    if (!!assignment && assignment !== deliveryId) {
      console.log("already assigned");
      return;
    }
    if (assignment === deliveryId) {
      console.log(`remove order ${id} from delivery ${deliveryId}`);
      markers.unassignOrder(id);
      // delete from orders list
      orders.splice(
        orders.findIndex((o) => o.id === id),
        1,
      );
      deselected_orders.add(id);
      changed = true;
      return;
    }
    console.log(`add order ${id} to delivery ${deliveryId}`);
    const new_order = markers.addOrderToDelivery(id, deliveryId);
    if (!new_order) return;
    orders.push(new_order.order);
    deselected_orders.delete(id);
    changed = true;
  });

  let lastsuccess = $state(0);
  $effect(() => {
    if (query.isPending) return;
    if (!query.isSuccess || lastsuccess >= query.dataUpdatedAt) return;
    lastsuccess = query.dataUpdatedAt;
    orders = query.data.data.orders;
  });
</script>

<section
  class="absolute top-[calc(100vh-24rem)] h-96 right-0 w-full bg-white p-3"
>
  <div class="flex justify-between mb-5">
    <div class="flex gap-2 items-center">
      <div
        class="aspect-square w-6 h-6 rounded-full border border-neutral-200"
        style:background-color={delivery?.color}
      ></div>
      <h1 class="font-bold">
        Liefertour {params.deliveryId}
        {delivery?.label}
      </h1>
      <div
        class="rounded-full border border-gray-600 px-1.5 flex items-center gap-1 h-fit"
      >
        <Clock size={16} class="text-gray-600" />{delivery?.until}
      </div>
      <button
        type="button"
        class="flex items-center"
        onclick={() => {
          open_edit_dialog = true;
        }}><Pencil size={16} /></button
      >
    </div>
    <button
      type="button"
      onclick={() => navigate("/deliveries", { viewTransition: true })}
      ><X /></button
    >
  </div>
  {#if changed}
    <button
      type="button"
      disabled={assignOrdersMutation.isPending}
      class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-300"
      onclick={() => {
        console.log("start saving");
        assignOrdersMutation.mutate({
          delivery_id: deliveryId,
          selected_orders: orders.map((o) => o.id),
          deselected_orders: Array.from(deselected_orders),
        });
      }}>Speichern</button
    >
  {/if}
  {#if !orders}
    Lädt...
  {/if}
  <table class="w-full">
    <thead class="bg-neutral-300 text-sm text-slate-700">
      <tr class="whitespace-nowrap text-left">
        <th>Id</th>
        <th>Lieferzeit</th>
        <th>#S</th>
        <th>#M</th>
        <th>#L</th>
        <th>Addresse</th>
        <th>PLZ</th>
        <th>Stadt</th>
        <th>Name</th>
      </tr>
    </thead>
    <tbody>
      {#each orders as order}
        <tr class="whitespace-nowrap odd:bg-blue-100">
          <td class="block max-w-16 overflow-hidden text-ellipsis px-2 py-1"
            >{order.id}</td
          >
          <td>{order.deliverymethod}</td>
          <td>{order.amountsizeS}</td>
          <td>{order.amountsizeM}</td>
          <td>{order.amountsizeL}</td>
          <td>{order.address}</td>
          <td>{order.postalcode}</td>
          <td>{order.city}</td>
          <td>{order.firstname} {order.lastname}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</section>

{#if open_edit_dialog}
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div
    class="absolute inset-0 bg-black/10"
    onclick={() => {
      open_edit_dialog = false;
    }}
  ></div>
  <div
    class="absolute inset-0 flex items-center justify-center pointer-events-none"
  >
    <div class="bg-white rounded shadow pointer-events-auto">
      <div class="flex justify-between p-2">
        <h1 class="font-bold">Liefertour {params.deliveryId} bearbeiten</h1>
        <button
          type="button"
          onclick={() => {
            open_edit_dialog = false;
          }}><X /></button
        >
      </div>
      {#if delivery}
        <EditDeliveryForm
          {delivery}
          on_saved={() => {
            open_edit_dialog = false;
          }}
        />
      {/if}
    </div>
  </div>
{/if}
