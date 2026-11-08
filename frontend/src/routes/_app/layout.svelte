<script lang="ts">
	import { p } from "sv-router/generated";
	import type { Snippet } from "svelte";
	import LogoutButton from "@/lib/auth/LogoutButton.svelte";
	import { setMarkerStoreContext } from "@/lib/map/markerstore.svelte";
	import MapComp from "@/lib/map/Map.svelte";
	import Markers from "@/lib/map/Markers.svelte";

	let { children }: { children: Snippet } = $props();

	setMarkerStoreContext();
</script>

<nav class="w-full">
	<div
		class="flex mx-auto max-w-2xl w-full gap-10 p-3 items-center justify-between"
	>
		<a href={p("/")} data-view-transition>Karte</a>
		<a href={p("/orders")} data-view-transition>Bestellungen</a>
		<a href={p("/deliveries")} data-view-transition>Liefertouren</a>
		<LogoutButton />
	</div>
</nav>
<div class="relative">
	<MapComp>
		<Markers />
	</MapComp>
	{@render children()}
</div>
