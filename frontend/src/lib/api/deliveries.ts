import { apiV1client } from "./axios";

export type Delivery = {
	id: number;
	label: string;
	until: string;
	color: string;
};

export type GetDeliveriesRes = {
	ok: boolean;
	error?: string;
	deliveries: Delivery[];
}

export type CreateDeliveryRes = {
	ok: boolean;
	error?: string;
	delivery: Delivery;
}

export type GetDeliveryRes = {
	ok: boolean;
	error?: string;
	delivery: Delivery;
}


export function getDeliveries() {
	return apiV1client.get<GetDeliveriesRes>("/api/v1/deliveries/all");
}

export function getDelivery(id: string) {
	return apiV1client.get<GetDeliveryRes>("/api/v1/deliveries", {
		params: {
			id
		}
	});
}
export function createNewDelivery(delivery: FormData) {
	return apiV1client.post<CreateDeliveryRes>("/api/v1/deliveries", delivery);
}
