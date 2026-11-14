<script lang="ts">
  import ColorPicker from "svelte-awesome-color-picker";
  import { generateRandomNonYellowHue, hslToRgb, rgbToHex } from "./colors";
  import { createMutation, useQueryClient } from "@tanstack/svelte-query";
  import { api } from "@/lib/api/api";
  import type { Delivery } from "../api/deliveries";

  interface Props {
    delivery: Delivery;
    on_saved: () => void;
  }
  const { delivery, on_saved }: Props = $props();

  function generateRandomColor() {
    const hue = generateRandomNonYellowHue();
    const sat = 0.7;
    const light = Math.random() * 0.2 + 0.5;
    const [r, g, b] = hslToRgb(hue, sat, light);
    const hex = rgbToHex(r, g, b);
    return hex;
  }

  function generateSwatches(n: number) {
    return Array.from({ length: n }).map(() => generateRandomColor());
  }

  const N_SWATCHES = 7;

  // svelte-ignore state_referenced_locally
  let hex = $state(delivery.color);
  let swatches = $state(generateSwatches(N_SWATCHES));

  const queryClient = useQueryClient();
  const editDeliveryMutation = createMutation(() => ({
    mutationFn: api.deliveries.editDelivery,
    onSuccess() {
      queryClient.invalidateQueries({
        queryKey: ["delivery", delivery.id.toString()],
      });
      on_saved();
    },
  }));
</script>

<form
  class="p-3 rounded flex flex-col w-fit"
  onsubmit={(e) => {
    e.preventDefault();
    const data = new FormData(e.currentTarget);
    data.append("id", delivery.id.toString());
    console.log(data);
    editDeliveryMutation.mutate(data);
  }}
>
  <div class="flex justify-between gap-5 mb-5">
    <div class="flex flex-col w-fit">
      <label for="color" class="font-bold text-sm">Farbe</label>
      <ColorPicker
        isAlpha={false}
        isTextInput={false}
        {swatches}
        label="Farbe wählen"
        bind:hex
        name="color"
      />
      <button
        class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow"
        type="button"
        onclick={() => {
          hex = generateRandomColor();
          swatches = generateSwatches(N_SWATCHES);
        }}>Mit Farben gambeln</button
      >
    </div>
    <div class="flex flex-col">
      <label for="until" class="font-bold text-sm">Lieferzeit</label>
      <select
        id="until"
        name="until"
        class="border border-gray-200 p-1 h-8 rounded"
        defaultvalue={delivery.until}
      >
        <option>bis 8 Uhr</option>
        <option>bis 10 Uhr</option>
      </select>
    </div>
    <div class="flex flex-col">
      <label for="label" class="font-bold text-sm">Notiz</label>
      <input
        id="label"
        name="label"
        type="text"
        class="border rounded py-1 px-2 border-gray-200"
        placeholder=" - "
        defaultvalue={delivery.label}
      />
    </div>
  </div>
  <button
    class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow"
    >Speichern</button
  >
</form>
