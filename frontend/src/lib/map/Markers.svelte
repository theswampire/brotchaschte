<script lang="ts">
  import { createQuery } from "@tanstack/svelte-query";
  import { api } from "@/lib/api/api";
  import { Marker } from "svelte-maplibre-gl";
  import { toLngLat, type Coords } from "@/lib/map/coords";
  import {
    getMarkerStoreContext,
    type MarkerEntry,
  } from "./markerstore.svelte";
  import MarkerSVG from "./MarkerSVG.svelte";
  import { STANDARD_COLORS } from "./marker";
  import Clock from "../icons/movingicons/clock.svelte";

  const markerstore = getMarkerStoreContext();

  const query = createQuery(() => ({
    queryKey: ["sorted-orders"],
    queryFn: api.orders.getSortedOrders,
  }));
  let lastsuccess = $state(0);

  $effect(() => {
    if (query.isPending) return;
    if (!query.isSuccess || lastsuccess >= query.dataUpdatedAt) return;
    lastsuccess = query.dataUpdatedAt;
    markerstore.setOrders(query.data.data.orders, query.data.data.deliveries);
  });

  let manual_order_coords = $derived.by<Coords | null>(() => {
    if (!markerstore.manual_order) return null;
    if (markerstore.manual_order.x.ok) {
      return [
        markerstore.manual_order.y.value,
        markerstore.manual_order.x.value,
      ];
    }
    if (markerstore.manual_order.suggested_easting.ok) {
      return [
        markerstore.manual_order.suggested_easting.value,
        markerstore.manual_order.suggested_northing.value,
      ];
    }
    return null;
  });

  let hoverstate = $state<{
    show: boolean;
    marker: MarkerEntry | null;
    coords: Coords;
  }>({
    show: false,
    marker: null,
    coords: [200, 200],
  });
</script>

{#each markerstore.markerlist.values() as m}
  <Marker
    offset={[0, -24 /* so it doesn't cover the address */]}
    lnglat={toLngLat([m.order.y, m.order.x])}
    onclick={() => {
      markerstore._markerOnClick(m.order.id);
    }}
  >
    {#snippet content()}
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        onmouseenter={(e) => {
          hoverstate.marker = m;
          hoverstate.coords = [e.pageX, e.pageY];
          hoverstate.show = true;
        }}
        onmouseleave={() => {
          hoverstate.show = false;
        }}
        style:transform-origin="bottom center"
        style:transform="rotate({m.rotation ?? 0}deg)"
      >
        <MarkerSVG color={m.color} hidden={!m.visible} />
      </div>
    {/snippet}
  </Marker>
{/each}

{#if manual_order_coords}
  <Marker
    offset={[0, -15 /* so it doesn't cover the address and is aligned */]}
    lnglat={toLngLat(manual_order_coords)}
    onclick={() => {
      markerstore._markerOnClick(markerstore.manual_order?.id as string);
    }}
  >
    {#snippet content()}
      <MarkerSVG color={STANDARD_COLORS.unassigned} />
    {/snippet}
  </Marker>
{/if}

{#if hoverstate.show}
  <div
    class="absolute bg-white p-4 z-10 -translate-x-1/2 -translate-y-62 rounded shadow border border-gray-200 text-base flex flex-col"
    style:top={`${hoverstate.coords[1]}px`}
    style:left={`${hoverstate.coords[0]}px`}
  >
    {#if hoverstate.marker?.delivery_id}
      <span class="font-bold text-xs"
        >Liefertour {hoverstate.marker?.delivery_id}</span
      >
    {:else}
      <span class="font-bold text-xs">keine Liefertour</span>
    {/if}

    <div class="flex gap-3 items-center justify-between">
      <div
        class="rounded-full border border-gray-600 px-1.5 flex items-center gap-1 h-fit w-fit data-[time=early]:border-red-500 data-[time=early]:text-red-600"
        data-time={hoverstate.marker?.order.deliverymethod === "8 Uhr"
          ? "early"
          : "late"}
      >
        <Clock size={16} />{hoverstate.marker?.order.deliverymethod}
      </div>
      <div class="grid grid-cols-3 gap-x-1">
        <span class="font-bold text-xs text-gray-500">#S</span>
        <span class="font-bold text-xs text-gray-500">#M</span>
        <span class="font-bold text-xs text-gray-500">#L</span>
        <span>{hoverstate.marker?.order.amountsizeS}</span>
        <span>{hoverstate.marker?.order.amountsizeM}</span>
        <span>{hoverstate.marker?.order.amountsizeL}</span>
      </div>
    </div>
    <span class="">
      {hoverstate.marker?.order.firstname}
      {hoverstate.marker?.order.lastname}
    </span>
    <span class="italic">
      {hoverstate.marker?.order.address}
    </span>
    <span class="italic">
      {hoverstate.marker?.order.postalcode}
      {hoverstate.marker?.order.city}
    </span>
  </div>
{/if}
