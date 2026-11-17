<script lang="ts">
	import { getMarkerStoreContext } from "@/lib/map/markerstore.svelte";

	const markers = getMarkerStoreContext();

	markers.setToDeliveryOverviewMode();
	markers.resetCallbacks();
	markers.restoreChanges();
</script>

