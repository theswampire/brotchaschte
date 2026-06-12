<script lang="ts">
	import { api } from "@/lib/api/api";
	import { createQuery } from "@tanstack/svelte-query";

	console.log("kek page");
	const query = createQuery(() => ({
		queryKey: ["kek"],
		queryFn: api.kek.getKek,
	}));
</script>

<section>
	Hehehehe

	{#if query.isPending}
		<p>Loading</p>
	{/if}
	{#if query.isSuccess}
		<p>{query.data}</p>
	{/if}
</section>
