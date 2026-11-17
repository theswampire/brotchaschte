<script lang="ts">
	import { getMarkerStoreContext } from "@/lib/map/markerstore.svelte";

	const markers = getMarkerStoreContext();

	markers.setToDeliveryOverviewMode();
	markers.resetCallbacks();
	markers.restoreChanges();
</script>

<section
	class="absolute inset-0 pointer-events-none h-full flex flex-col justify-center"
>
</section>
