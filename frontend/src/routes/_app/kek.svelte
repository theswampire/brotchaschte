<script lang="ts">
	import { createQuery } from "@tanstack/svelte-query";
	import { api } from "@/lib/api/api";

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
