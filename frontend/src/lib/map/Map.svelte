<script lang="ts">
    import "svelte-maplibre-gl/vite";
    import { MapLibre, Marker } from "svelte-maplibre-gl";
    import { type Coords, toLngLat } from "@/lib/map/coords";
    import { STANDARD_COLORS } from "@/lib/map/marker";
    import type { Map, StyleSpecification } from "maplibre-gl";

    interface Props {}

    const MAPLIBRE_URL =
        "https://vectortiles.geo.admin.ch/styles/ch.swisstopo.basemap.vt/style.json";

    async function getStyle() {
        const res = await fetch(MAPLIBRE_URL);
        const style = (await res.json()) as StyleSpecification;
        const hiddenLayerIds = [
            "public_transport",
            "poi_rank2",
            "poi_rank1",
            "poi_rank3_rotation_flat",
            "poi_motorway",
            "road_number",
            "aerodrome_label",
            "peak_rank2",
        ];
        style.layers = style.layers.filter(
            (layer) => !hiddenLayerIds.includes(layer.id),
        );

        return style;
    }
    const style = await getStyle();
    let map: Map | undefined = $state(undefined);
    const MAT_COORDS: Coords = [2699530.88, 1262572.69];
</script>

<MapLibre
    bind:map
    center={toLngLat(MAT_COORDS)}
    class="h-[calc(100vh-6em)]"
    {style}
    zoom={14}
>
    <Marker lnglat={toLngLat(MAT_COORDS)} color={STANDARD_COLORS.center}>
        {#snippet content()}
            <div
                class="bg-white size-5 rounded-full absolute top-2 left-1/2 -translate-x-1/2"
            ></div>
            <div
                class="bg-black size-2 scale-y-[0.5] blur-sm rounded-full absolute bottom-0 left-1/2 -translate-x-1/2"
            ></div>
            <div
                class="bg-black size-1 scale-y-[0.5] blur-sm rounded-full absolute bottom-0 left-1/2 -translate-x-1/2"
            ></div>
            <div
                class="bg-black size-1 scale-y-[0.5] blur-sm rounded-full absolute bottom-0 left-1/2 -translate-x-1/2"
            ></div>
            <div
                class="bg-black size-1 scale-y-[0.5] blur-sm rounded-full absolute bottom-0 left-1/2 -translate-x-1/2"
            ></div>
            <div
                class="bg-black/20 size-2 scale-y-[0.5] rounded-full absolute bottom-0 left-1/2 -translate-x-1/2"
            ></div>
            <img
                src="/assets/images/home.png"
                alt="Home Marker"
                class="size-11 relative"
            />
        {/snippet}
    </Marker>
</MapLibre>
