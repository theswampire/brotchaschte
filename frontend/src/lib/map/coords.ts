import proj4 from "proj4";

// register LV95 (EPSG:2056) coordinate system
proj4.defs(
	"EPSG:2056",
	"+proj=somerc +lat_0=46.95240555555556 +lon_0=7.439583333333333 +k_0=1 +x_0=2600000 +y_0=1200000 +ellps=bessel +towgs84=674.374,15.056,405.346,0,0,0,0 +units=m +no_defs",
);

export type Coords = [number, number];

/** Converts Swiss LV95 coords like (2'698'715.67, 1'262'658.47) to Longiture Latitude like (8.74, 47.5) */
export function toLngLat(coords: Coords): Coords {
	return proj4("EPSG:2056", "EPSG:4326", coords);
}

/** Converts Longiture Latitude like (8.74, 47.5) to Swiss LV95 coords like (2'698'715.67, 1'262'658.47)  */
export function toCH(coords: Coords): Coords {
	return proj4("EPSG:4326", "EPSG:2056", coords);
}
