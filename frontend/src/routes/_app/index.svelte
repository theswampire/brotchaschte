<script lang="ts">
	import MapComp from "@/lib/map/Map.svelte";
    import GeocodeOrdersButton from "@/lib/orders/GeocodeOrdersButton.svelte";
	import OrdersTable from "@/lib/orders/OrdersTable.svelte";
    import SyncButton from "@/lib/orders/SyncButton.svelte";

	let visible_orders = $state(false);
</script>

<section>
	<button
		type="button"
		class="shadow p-1 px-2 rounded border border-gray-300 m-2 w-48"
		onclick={() => {
			visible_orders = !visible_orders;
		}}
	>
		{#if visible_orders}
			Karte zeigen
		{/if}
		{#if !visible_orders}
			Bestellungen zeigen
		{/if}
	</button>
	<div class="relative">
		<MapComp />
		<section
			data-visible={visible_orders.valueOf()}
			class="bg-white absolute top-0 p-3 w-full h-full data-[visible=false]:opacity-0 data-[visible=false]:pointer-events-none mx-auto transition-opacity overflow-hidden"
		>
			<div class="flex justify-between">
				<h1 class="font-bold text-3xl">Bestellungen</h1>
				<GeocodeOrdersButton />
				<SyncButton />
				<button type="button" onclick={() => (visible_orders = false)}
					>X</button
				>
			</div>
			<div class="pb-5 overflow-y-auto h-full">
				<OrdersTable />
			</div>
		</section>
	</div>
</section>
