<script lang="ts">
	import GenerateOrdersButton from "@/lib/orders/GenerateOrdersButton.svelte";
	import GeocodeOrdersButton from "@/lib/orders/GeocodeOrdersButton.svelte";
	import OrdersTable from "@/lib/orders/OrdersTable.svelte";
	import SyncButton from "@/lib/orders/SyncButton.svelte";
</script>

<section
	class="bg-white absolute top-0 p-3 w-full h-full mx-auto transition-opacity overflow-hidden"
>
	<div class="flex justify-between">
		<h1 class="font-bold text-3xl">Bestellungen</h1>
		<GenerateOrdersButton />
		<GeocodeOrdersButton />
		<SyncButton />
	</div>
	<div class="pb-5 overflow-y-auto h-full">
		<OrdersTable />
	</div>
</section>
