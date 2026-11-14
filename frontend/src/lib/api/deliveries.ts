import { apiV1client } from "./axios";
import type { Order } from "./orders";

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
};

export type CreateDeliveryRes = {
  ok: boolean;
  error?: string;
  delivery: Delivery;
};

export type GetDeliveryRes = {
  ok: boolean;
  error?: string;
  delivery: Delivery;
  orders: Order[];
};
export type AssignOrdersRes = {
  ok: boolean;
  error?: string;
};
export type EditDeliveryRes = {
  ok: boolean;
  error?: string;
  delivery: Delivery;
};

export function getDeliveries() {
  return apiV1client.get<GetDeliveriesRes>("/api/v1/deliveries/all");
}

export function getDelivery(id: string) {
  return apiV1client.get<GetDeliveryRes>("/api/v1/deliveries", {
    params: {
      id,
    },
  });
}
export function createNewDelivery(delivery: FormData) {
  return apiV1client.post<CreateDeliveryRes>("/api/v1/deliveries", delivery);
}
export function assignOrdersToDelivery(data: {
  delivery_id: number;
  selected_orders: string[];
  deselected_orders: string[];
}) {
  return apiV1client.post<AssignOrdersRes>("/api/v1/deliveries/assign", {
    delivery_id: data.delivery_id,
    len_selected: data.selected_orders.length,
    len_deselected: data.deselected_orders.length,
    selected_order_ids: data.selected_orders,
    deselected_order_ids: data.deselected_orders,
  });
}

export function editDelivery(data: FormData) {
  return apiV1client.patch<EditDeliveryRes>("/api/v1/deliveries", data);
}
