<script lang="ts">
  import { p } from "sv-router/generated";
  import type { Snippet } from "svelte";
  import LogoutButton from "@/lib/auth/LogoutButton.svelte";
  import { setMarkerStoreContext } from "@/lib/map/markerstore.svelte";
  import MapComp from "@/lib/map/Map.svelte";
  import Markers from "@/lib/map/Markers.svelte";
  import type { Map } from "maplibre-gl";
  import { setMapHandleContext } from "@/lib/map/mapcontext";

  let { children }: { children: Snippet } = $props();

  let map: { instance?: Map } = $state({});

  setMarkerStoreContext();
  setMapHandleContext(map);
</script>

<nav class="w-full h-14 border-b border-gray-200">
  <div
    class="flex mx-auto max-w-2xl h-full w-full gap-10 items-center justify-between"
  >
    <a href={p("/")} data-view-transition>Karte</a>
    <a href={p("/orders")} data-view-transition>Bestellungen</a>
    <a href={p("/deliveries")} data-view-transition>Liefertouren</a>
    <LogoutButton />
  </div>
</nav>
<!-- h-[calc(100%-56px)] is the exact height left after the nav element with h-14 -->
<div class="relative h-[calc(100%-56px)]">
  <MapComp bind:map={map.instance}>
    <Markers />
  </MapComp>
  {@render children()}
</div>
