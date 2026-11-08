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


export function getDeliveries() {
	return apiV1client.get<GetDeliveriesRes>("/api/v1/deliveries");
}
