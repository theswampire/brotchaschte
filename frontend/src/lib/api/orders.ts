import { apiV1client } from "./axios";

export type Order = {
	id: string
	lastname: string
	firstname: string
	phonenumber: string
	email: string
	amountsizeS: number
	amountsizeM: number
	amountsizeL: number
	address: string
	city: string
	postalcode: string
	deliverymethod: string
	remarks: string
	ordermethod: string
	vendorgroup: string
}

export type GetOrdersRes = {
	ok: boolean;
	orders?: Order[];
	error?: string;
};

export function getOrders() {
	return apiV1client.get<GetOrdersRes>("/api/v1/orders");
}

export function generateRandomOrders() {
	return apiV1client.post<GetOrdersRes>("/api/v1/orders/generate")
}

