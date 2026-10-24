import { apiV1client } from "./axios";

export type Order = {
	id: string;
	lastname: string;
	firstname: string;
	phonenumber: string;
	email: string;
	amountsizeS: number;
	amountsizeM: number;
	amountsizeL: number;
	address: string;
	city: string;
	postalcode: string;
	deliverymethod: string;
	remarks: string;
	ordermethod: string;
	vendorgroup: string;
};

export type FailedGeocodedOrder = Order & {
	message: string;
	suggested_northing: number | null;
	suggested_easting: number | null;
	suggested_label: string | null;
};
export type GeocodedOrder = Order & {
	x: number;
	y: number;
	label: string;
};

export type GetOrdersRes = {
	ok: boolean;
	not_geocoded_orders?: Order[];
	failed_geocoded_orders?: FailedGeocodedOrder[];
	geocoded_orders?: GeocodedOrder[];
	error?: string;
};

export type GeocodeOrdersRes = {
	ok: boolean;
	error?: string;
};

export function getOrders() {
	return apiV1client.get<GetOrdersRes>("/api/v1/orders").then(x => {console.log(x); return x});
}

export function generateRandomOrders() {
	return apiV1client.post<GetOrdersRes>("/api/v1/orders/generate");
}

export function pullUpstreamOrders() {
	return apiV1client.post<GetOrdersRes>("/api/v1/orders/sync");
}

export function geocodeOrders() {
	return apiV1client.post<GeocodedOrder>("/api/v1/orders/geolocate");
}
