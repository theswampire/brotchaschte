<script lang="ts">
  import { api } from "@/lib/api/api";
  import { type GeocodedWithSuggestionsOrder } from "@/lib/api/orders";
  import X from "@/lib/icons/movingicons/x.svelte";
  import { toCH, toLngLat, type Coords } from "@/lib/map/coords";
  import { getMapHandleContext } from "@/lib/map/mapcontext";
  import { getMarkerStoreContext } from "@/lib/map/markerstore.svelte";
  import {
    createMutation,
    createQuery,
    useQueryClient,
  } from "@tanstack/svelte-query";
  import type { MapLayerEventType } from "maplibre-gl";
  import { navigate, route } from "sv-router/generated";

  const params = route.getParams("/orders/:orderId/manual");
  const query = createQuery(() => ({
    queryKey: ["order", params.orderId],
    queryFn: api.orders.getGeocodedWithSuggestionsOrder.bind(
      null,
      params.orderId,
    ),
  }));
  const markers = getMarkerStoreContext();
  const map = getMapHandleContext();
  markers.setToManualLocateMode();

  let changed = $state(false);
  let lastsuccess = $state(0);

  // select coords if existing with non-suggestion preferred
  function getCoords(order: GeocodedWithSuggestionsOrder): Coords {
    if (order.x.ok) {
      return [order.y.value, order.x.value];
    }
    return [order.suggested_easting.value, order.suggested_northing.value];
  }

  $effect(() => {
    if (query.isPending) return;
    if (!query.isSuccess || lastsuccess >= query.dataUpdatedAt) return;
    lastsuccess = query.dataUpdatedAt;

    markers.addManualLocateOrder(query.data.data.order);

    const coords: Coords = getCoords(query.data.data.order);
    if (
      query.data.data.order.suggested_northing.ok ||
      query.data.data.order.x.ok
    ) {
      map.instance?.flyTo({
        center: toLngLat(coords),
        zoom: 18,
      });
    }
  });

  $effect(() => {
    if (!map.instance) return;

    const handler = (ev: MapLayerEventType["click"]) => {
      if (!query.data) {
        console.warn("tried setting manual coordinates while query is empty");
        return;
      }
      const coords = toCH([ev.lngLat.lng, ev.lngLat.lat]);
      markers.addManualLocateOrder({
        ...query.data!.data.order,
        y: { ok: true, value: coords[0] },
        x: { ok: true, value: coords[1] },
        suggested_easting: { ok: true, value: coords[0] },
        suggested_northing: { ok: true, value: coords[1] },
        suggested_label: "",
      });
      changed = true;
    };

    map.instance.on("click", handler);

    return () => map.instance?.off("click", handler);
  });

  const queryClient = useQueryClient();
  let adoptUnchangedMutation = createMutation(() => ({
    mutationFn: api.orders.manuallyLocateOrder,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      queryClient.invalidateQueries({ queryKey: ["order", params.orderId] });

      navigate("/orders", { viewTransition: true });
    },
  }));
  let locateAndUpdateAddressMutation = createMutation(() => ({
    mutationFn: api.orders.manuallyLocateWithAddressUpdateOrder,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      queryClient.invalidateQueries({ queryKey: ["order", params.orderId] });
      navigate("/orders", { viewTransition: true });
    },
  }));

  function adoptUnchangedHandler() {
    const data = new FormData();
    if (!markers.manual_order) {
      console.warn("adopt unchanged for null manual_order");
      return;
    }
    if (
      !markers.manual_order.suggested_northing.ok ||
      !markers.manual_order.suggested_easting.ok
    ) {
      console.warn("adopt unchanged for unset suggested_northing");
      return;
    }
    const coords = getCoords(markers.manual_order);
    data.set("easting", coords[0].toString());
    data.set("northing", coords[1].toString());
    data.set("order_id", markers.manual_order.id);
    adoptUnchangedMutation.mutate(data);
  }

  function manualLocateHandler(
    ev: SubmitEvent & {
      currentTarget: EventTarget & HTMLFormElement;
    },
  ) {
    ev.preventDefault();

    const data = new FormData(ev.currentTarget);
    if (!markers.manual_order) {
      console.warn("manual locate for null manual_order");
      return;
    }
    if (
      !markers.manual_order.suggested_northing.ok ||
      !markers.manual_order.suggested_easting.ok
    ) {
      console.warn("manual locate for unset suggested_northing");
      return;
    }
    data.set(
      "northing",
      markers.manual_order.suggested_northing.value.toString(),
    );
    data.set(
      "easting",
      markers.manual_order.suggested_easting.value.toString(),
    );
    data.set("order_id", markers.manual_order.id);
    locateAndUpdateAddressMutation.mutate(data);
  }
</script>

<section
  class="absolute right-10 top-10 bg-white p-3 rounded border border-gray-200"
>
  <div class="w-full flex justify-between">
    <h1 class="font-bold text-lg">Bestellung manuell lokalisieren</h1>
    <button
      type="button"
      onclick={() => navigate("/orders", { viewTransition: true })}
      ><X /></button
    >
  </div>
  <span class="text-sm text-gray-400">{params.orderId}</span>

  {#if query.isSuccess}
    <span class="text-sm text-gray-400"
      >{query.data.data.order.firstname} {query.data.data.order.lastname}</span
    >
    <form class="flex flex-col" onsubmit={manualLocateHandler}>
      <label for="address" class="font-bold text-sm"
        >Strasse und Hausnummer</label
      >
      <input
        id="address"
        class="border rounded py-1 px-2 border-gray-200 mb-2"
        name="address"
        type="text"
        defaultvalue={query.data.data.order.address}
        onchange={() => (changed = true)}
      />
      <label for="postalcode" class="font-bold text-sm">Postleitzahl</label>
      <input
        id="postalcode"
        class="border rounded py-1 px-2 border-gray-200 mb-2"
        name="postalcode"
        type="text"
        defaultvalue={query.data.data.order.postalcode}
        onchange={() => (changed = true)}
      />
      <label for="city" class="font-bold text-sm">Ort</label>
      <input
        id="city"
        class="border rounded py-1 px-2 border-gray-200"
        name="city"
        type="text"
        defaultvalue={query.data.data.order.city}
        onchange={() => (changed = true)}
      />
      {#if query.data.data.order.suggested_label}
        <div class="mt-5">
          <p>
            Die angegebene Adresse wurde nicht gefunden. Hier ist ein Vorschlag
            stattdessen.
          </p>
          <span class="font-bold">Vorschlag:</span>
          <span class="italic text-gray-700">
            {query.data.data.order.suggested_label}
          </span>
        </div>
      {/if}

      <button
        disabled={!changed ||
          adoptUnchangedMutation.isPending ||
          locateAndUpdateAddressMutation.isPending}
        class="border mt-5 rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-200 disabled:text-gray-500 disabled:shadow-none"
        type="submit"
        >Speichern

        {#if changed}
          *
        {/if}
      </button>

      {#if query.data.data.order.suggested_label && !changed}
        <button
          disabled={adoptUnchangedMutation.isPending}
          type="button"
          onclick={adoptUnchangedHandler}
          class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-200 disabled:text-gray-500 disabled:shadow-none mt-5"
          >Übernehmen so wie es ist
        </button>
      {/if}
    </form>
  {/if}
</section>
