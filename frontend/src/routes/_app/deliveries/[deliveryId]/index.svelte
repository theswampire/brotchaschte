<script lang="ts">
	import { api } from "@/lib/api/api";
	import { createQuery } from "@tanstack/svelte-query";
	import { route } from "sv-router/generated";
	const params = route.getParams("/deliveries/:deliveryId");

	const query = createQuery(() => ({
		queryKey: ["delivery", params.deliveryId],
		queryFn: () => api.deliveries.getDelivery(params.deliveryId),
	}));
</script>

<section class="absolute top-[calc(100vh-24rem)] h-96 right-0 w-full bg-white">
	<h1>Delivery: {params.deliveryId}</h1>
	{#if query.isLoading}
		Lädt...
	{/if}
	{#if query.isSuccess}
		{query.data.data.delivery.label}
		{query.data.data.delivery.until}
	{/if}
</section>
