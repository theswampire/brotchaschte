<script lang="ts">
	import DeliveryOverview from "@/lib/deliveries/DeliveryOverview.svelte";
	import { getMarkerStoreContext } from "@/lib/map/markerstore.svelte";
	const markers = getMarkerStoreContext();
	markers.setToDeliveryOverviewMode();
	markers.resetCallbacks();
	markers.restoreChanges();
</script>

<DeliveryOverview />
