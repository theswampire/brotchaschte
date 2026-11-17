<script lang="ts">
  import type { Order } from "@/lib/api/orders";
  import Dialog from "@/lib/Dialog.svelte";
  import X from "@/lib/icons/movingicons/x.svelte";
  import EditOrdersForm from "@/lib/orders/EditOrdersForm.svelte";
  //import GenerateOrdersButton from "@/lib/orders/GenerateOrdersButton.svelte";
  import GeocodeOrdersButton from "@/lib/orders/GeocodeOrdersButton.svelte";
  import OrdersTable from "@/lib/orders/OrdersTable.svelte";
  import SyncButton from "@/lib/orders/SyncButton.svelte";

  let show_edit_order = $state(false);
  let current_edit_order = $state<Order | null>(null);

  function edit_order_callback(order: Order) {
    console.log(order);
    current_edit_order = order;
    show_edit_order = true;
  }
</script>

<Dialog bind:show={show_edit_order} class="mx-auto -translate-1/2 top-1/2">
  <div class="bg-white">
    <div class="flex justify-between p-3 items-start">
      <div>
        <h3 class="font-bold text-xl">Bestellung bearbeiten</h3>
        <span class="text-sm text-gray-400">{current_edit_order?.id}</span>
      </div>
      <button onclick={() => (show_edit_order = false)}><X /></button>
    </div>
    <EditOrdersForm
      order={current_edit_order}
      on_edited_order_callback={() => (show_edit_order = false)}
    />
  </div>
</Dialog>

<section
  class="bg-white absolute top-0 p-3 w-full h-full mx-auto transition-opacity overflow-hidden"
>
  <div class="flex justify-between">
    <h1 class="font-bold text-3xl">Bestellungen</h1>
    <!-- <GenerateOrdersButton /> -->
    <GeocodeOrdersButton />
    <SyncButton />
  </div>
  <div class="pb-5 overflow-y-auto h-full">
    <OrdersTable {edit_order_callback} />
  </div>
</section>
