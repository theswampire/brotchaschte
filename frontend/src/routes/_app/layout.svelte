<script lang="ts">
	import { p } from "sv-router/generated";
	import type { Snippet } from "svelte";
	import LogoutButton from "@/lib/auth/LogoutButton.svelte";

	let { children }: { children: Snippet } = $props();
</script>

<nav class="w-full">
	<div class="flex mx-auto max-w-2xl gap-10 p-3 items-center">
		<a href={p("/")}>Karte</a>
		<a href={p("/orders")}>Bestellungen</a>
		<LogoutButton />
	</div>
</nav>
<div class="wrapper">
	{@render children()}
</div>
