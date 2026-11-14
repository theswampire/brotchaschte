<script lang="ts">
  import { createQuery, useQueryClient } from "@tanstack/svelte-query";
  import { api } from "@/lib/api/api";
  import CreateDeliveryForm from "./CreateDeliveryForm.svelte";
  import { navigate } from "sv-router/generated";
  import X from "../icons/movingicons/x.svelte";
  import EditDeliveryForm from "./EditDeliveryForm.svelte";
  import type { Delivery } from "../api/deliveries";
  import Pencil from "../icons/movingicons/pencil.svelte";
  import Clock from "../icons/movingicons/clock.svelte";

  const query = createQuery(() => ({
    queryKey: ["deliveries"],
    queryFn: api.deliveries.getDeliveries,
  }));
  const queryClient = useQueryClient();

  let edit_dialog_state = $state<{ open: boolean; delivery?: Delivery }>({
    open: false,
  });
</script>

<section
  class="bg-white left-1/2 top-0 w-1/2 h-full shadow border-l-gray-300 absolute p-4"
>
  <h2 class="font-bold text-3xl">Lieferungen</h2>

  <CreateDeliveryForm />

  {#if query.isPending}
    <p>Lädt...</p>
  {/if}
  {#if query.isSuccess}
    <table class="w-full">
      <thead class="bg-neutral-300 text-sm text-slate-700">
        <tr class="whitespace-nowrap text-left">
          <th class="w-10"></th>
          <th>Id</th>
          <th>Notiz</th>
          <th>Lieferzeit</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each query.data.data.deliveries as delivery}
          <tr
            onclick={() =>
              navigate("/deliveries/:deliveryId", {
                params: { deliveryId: delivery.id.toString() },
                viewTransition: true,
              })}
            class="hover:bg-blue-50 transition-colors h-12 border-b"
          >
            <td class="w-fit p-0">
              <div
                class="aspect-square w-5 rounded-full mx-auto"
                style:background-color={delivery.color}
              ></div>
            </td>
            <td>
              {delivery.id}
            </td>
            <td>
              {delivery.label}
            </td>
            <td>
              <div
                class="rounded-full border border-gray-600 w-28 px-1.5 flex items-center justify-start gap-1 h-fit"
              >
                <Clock size={16} class="text-gray-600" />{delivery.until}
              </div>
            </td>
            <td>
              <button
                type="button"
                onclick={() => {
                  edit_dialog_state.open = true;
                  edit_dialog_state.delivery = delivery;
                }}><Pencil size={16} /></button
              >
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</section>

{#if edit_dialog_state.open}
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div
    class="absolute inset-0 bg-black/10"
    onclick={() => {
      edit_dialog_state.open = false;
    }}
  ></div>
  <div
    class="absolute inset-0 flex items-center justify-center pointer-events-none"
  >
    <div class="bg-white rounded shadow pointer-events-auto">
      <div class="flex justify-between p-2">
        <h1 class="font-bold">
          Liefertour {edit_dialog_state.delivery?.id} bearbeiten
        </h1>
        <button
          type="button"
          onclick={() => {
            edit_dialog_state.open = false;
          }}><X /></button
        >
      </div>
      {#if edit_dialog_state.delivery}
        <EditDeliveryForm
          delivery={edit_dialog_state.delivery}
          on_saved={() => {
            edit_dialog_state.open = false;
            queryClient.invalidateQueries({ queryKey: ["deliveries"] });
          }}
        />
      {/if}
    </div>
  </div>
{/if}
