<script lang="ts">
  import type { Order } from "@/lib/api/orders";
  import { createMutation, useQueryClient } from "@tanstack/svelte-query";
  import { api } from "@/lib/api/api";
  import { navigate } from "sv-router/generated";
  import MarkerSVG from "../map/MarkerSVG.svelte";

  interface Props {
    order: Order | null;
    on_edited_order_callback?: () => void;
  }
  const { order, on_edited_order_callback }: Props = $props();

  const queryClient = useQueryClient();
  const updateOrderMutation = createMutation(() => ({
    mutationFn: api.orders.updateOrder,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      on_edited_order_callback?.();
    },
  }));
</script>

<form
  class="flex flex-col p-3"
  onsubmit={(e) => {
    e.preventDefault();
    const data = new FormData(e.currentTarget);
    updateOrderMutation.mutate(data);
  }}
>
  <input name="id" value={order?.id} hidden />
  <div class="flex gap-2">
    <div>
      <label for="firstname" class="block">Vorname</label>
      <input
        id="firstname"
        type="text"
        name="firstname"
        class="border rounded py-1 px-2 border-gray-200"
        defaultvalue={order?.firstname}
      />
    </div>
    <div>
      <label for="lastname" class="block">Nachname</label>
      <input
        id="lastname"
        type="text"
        name="lastname"
        class="border rounded py-1 px-2 border-gray-200"
        defaultvalue={order?.lastname}
      />
    </div>
  </div>
  <div class="flex gap-2">
    <div>
      <label for="email" class="block">E-Mail</label>
      <input
        id="email"
        type="text"
        name="email"
        class="border rounded py-1 px-2 border-gray-200"
        defaultvalue={order?.email}
      />
    </div>
    <div>
      <label for="phone" class="block">Telefon</label>
      <input
        id="phone"
        type="text"
        name="phone"
        class="border rounded py-1 px-2 border-gray-200"
        defaultvalue={order?.phonenumber}
      />
    </div>
  </div>

  <label for="address">Adresse</label>
  <input
    id="address"
    type="text"
    name="address"
    class="border rounded py-1 px-2 border-gray-200"
    defaultvalue={order?.address}
  />

  <div class="flex gap-2">
    <div>
      <label for="postal" class="block">Postleitzahl</label>
      <input
        id="postal"
        type="text"
        name="postal"
        class="border rounded py-1 px-2 border-gray-200"
        defaultvalue={order?.postalcode}
      />
    </div>
    <div>
      <label for="city" class="block">Stadt</label>
      <input
        id="city"
        type="text"
        name="city"
        defaultvalue={order?.city}
        class="border rounded py-1 px-2 border-gray-200"
      />
    </div>
  </div>
  <div class="grid grid-cols-3">
    <label for="amount_s">#S</label>
    <label for="amount_m">#M</label>
    <label for="amount_l">#L</label>
    <input
      id="amount_s"
      type="number"
      name="amount_s"
      class="border rounded py-1 px-2 border-gray-200 w-20"
      defaultvalue={order?.amountsizeS}
    />
    <input
      id="amount_m"
      type="number"
      name="amount_m"
      class="border rounded py-1 px-2 border-gray-200 w-20"
      defaultvalue={order?.amountsizeM}
    />
    <input
      id="amount_l"
      type="number"
      name="amount_l"
      class="border rounded py-1 px-2 border-gray-200 w-20"
      defaultvalue={order?.amountsizeL}
    />
  </div>
  <label for="delivery_method">Lieferzeit</label>
  <select
    defaultValue={order?.deliverymethod}
    id="delivery_method"
    name="delivery_method"
    class="border rounded py-1 px-2 border-gray-200"
  >
    <option value="8 Uhr">8 Uhr</option>
    <option value="10 Uhr">10 Uhr</option>
  </select>
  <label for="remarks">Bemerkungen</label>
  <textarea
    id="remarks"
    name="remarks"
    class="border rounded py-1 px-2 border-gray-200"
    defaultvalue={order?.remarks}></textarea>

  <button
    type="submit"
    disabled={updateOrderMutation.isPending}
    class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-200 disabled:text-gray-500 disabled:shadow-none mt-5"
    >Speichern</button
  >
  <button
    type="button"
    disabled={updateOrderMutation.isPending}
    class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-200 disabled:text-gray-500 disabled:shadow-none mt-5"
    onclick={() => {
      if (!order) return;
      navigate("/orders/:orderId/manual", {
        params: { orderId: order.id },
        viewTransition: true,
      });
    }}
  >
    <div class="inline-block h-4 scale-50 -translate-y-1.5">
      <MarkerSVG color="black" />
    </div>
    Lokalisierung</button
  >
</form>
