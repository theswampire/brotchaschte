<script lang="ts">
  import { createQuery } from "@tanstack/svelte-query";
  import { api } from "@/lib/api/api";
  import { Marker } from "svelte-maplibre-gl";
  import { toLngLat, type Coords } from "@/lib/map/coords";
  import { getMarkerStoreContext } from "./markerstore.svelte";
  import MarkerSVG from "./MarkerSVG.svelte";
  import { STANDARD_COLORS } from "./marker";

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
</script>

{#each markerstore.markerlist.values() as m}
  {#if m.visible}
    <Marker
      offset={[0, -24 /* so it doesn't cover the address */]}
      lnglat={toLngLat([m.order.y, m.order.x])}
      onclick={() => {
        markerstore._markerOnClick(m.order.id);
      }}
    >
      {#snippet content()}
        <MarkerSVG color={m.color} />
      {/snippet}
    </Marker>
  {/if}
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
