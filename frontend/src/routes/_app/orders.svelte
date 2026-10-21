<script lang="ts">
	import { createQuery } from "@tanstack/svelte-query";
	import { api } from "@/lib/api/api";
	import GenerateOrdersButton from "@/lib/orders/GenerateOrdersButton.svelte";

	const query = createQuery(() => ({
		queryKey: ["orders"],
		queryFn: api.orders.getOrders,
	}));
</script>

<section>
	<h1>Bestellungen</h1>
	<GenerateOrdersButton />

	{#if query.isPending}
		<p>Loading</p>
	{/if}
	{#if query.isSuccess}
		<div class="w-full overflow-x-auto">
			<table class="w-full">
				<thead class="bg-neutral-300 text-sm text-slate-700">
					<tr class="whitespace-nowrap text-left">
						<th>Id</th>
						<th>Name</th>
						<th>E-Mail</th>
						<th>Tel</th>
						<th>Addresse</th>
						<th>PLZ</th>
						<th>Stadt</th>
						<th>#S</th>
						<th>#M</th>
						<th>#L</th>
					</tr>
				</thead>
				<tbody>
					{#each query.data.data.orders as o}
						<tr class="whitespace-nowrap even:bg-blue-200">
							<td class="block max-w-20 overflow-hidden text-ellipsis px-2 py-1">{o.id}</td>
							<td>{o.firstname} {o.lastname}</td>
							<td>{o.email}</td>
							<td>{o.phonenumber}</td>
							<td>{o.address}</td>
							<td>{o.postalcode}</td>
							<td>{o.city}</td>
							<td>{o.amountsizeS}</td>
							<td>{o.amountsizeM}</td>
							<td>{o.amountsizeL}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</section>
