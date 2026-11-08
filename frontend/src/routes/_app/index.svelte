<script lang="ts">
	import DeliveryOverview from "@/lib/deliveries/DeliveryOverview.svelte";
	import MapComp from "@/lib/map/Map.svelte";
	import Markers from "@/lib/map/Markers.svelte";
	import { setMarkerStoreContext } from "@/lib/map/markerstore.svelte";
	import GenerateOrdersButton from "@/lib/orders/GenerateOrdersButton.svelte";
	import GeocodeOrdersButton from "@/lib/orders/GeocodeOrdersButton.svelte";
	import OrdersTable from "@/lib/orders/OrdersTable.svelte";
	import SyncButton from "@/lib/orders/SyncButton.svelte";

	let tab = $state<"orders" | "deliveries" | "map">("map");
	const markerstore = setMarkerStoreContext();
</script>

<section>
	<button
		type="button"
		class="shadow p-1 px-2 rounded border border-gray-300 m-2 w-48"
		onclick={() => {
			tab = "orders";
		}}
	>
		Bestellungen zeigen
	</button>
	<button
		type="button"
		class="shadow p-1 px-2 rounded border border-gray-300 m-2 w-48"
		onclick={() => {
			tab = "deliveries";
		}}
	>
		Lieferungen zeigen
	</button>
	<button
		type="button"
		class="shadow p-1 px-2 rounded border border-gray-300 m-2 w-48"
		onclick={() => {
			tab = "map";
		}}
	>
		Karte zeigen
	</button>
	<button
		type="button"
		class="shadow p-1 px-2 rounded border border-gray-300 m-2 w-48"
		onclick={() => markerstore.setToDeliveryOverviewMode()}
		>Switch to group</button
	>
	<button
		type="button"
		class="shadow p-1 px-2 rounded border border-gray-300 m-2 w-48"
		onclick={() => markerstore.setToDeliveryEditMode(12)}
		>Switch to delivery</button
	>
	<div class="relative">
		<MapComp>
			<Markers />
		</MapComp>
		<section
			data-visible={tab === "orders"}
			class="bg-white absolute top-0 p-3 w-full h-full data-[visible=false]:opacity-0 data-[visible=false]:pointer-events-none mx-auto transition-opacity overflow-hidden"
		>
			<div class="flex justify-between">
				<h1 class="font-bold text-3xl">Bestellungen</h1>
				<GenerateOrdersButton />
				<GeocodeOrdersButton />
				<SyncButton />
				<button type="button" onclick={() => (tab = "map")}>X</button>
			</div>
			<div class="pb-5 overflow-y-auto h-full">
				<OrdersTable />
			</div>
		</section>

		<DeliveryOverview visible={tab === "deliveries"} />
	</div>
</section>
