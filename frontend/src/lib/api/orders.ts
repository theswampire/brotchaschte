import { apiV1client } from "./axios";
import type { Delivery } from "./deliveries";

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

export type MaybeNumber = {
  ok: boolean;
  value: number;
};

export type FailedGeocodedOrder = Order & {
  message: string;
  suggested_northing: MaybeNumber;
  suggested_easting: MaybeNumber;
  suggested_label: string | null;
};

export type GeocodedOrder = Order & {
  x: number;
  y: number;
  label: string;
};

export type GeocodedWithSuggestionsOrder = FailedGeocodedOrder & {
  x: MaybeNumber;
  y: MaybeNumber;
  label: string;
};

export type GetOrdersRes = {
  ok: boolean;
  not_geocoded_orders?: Order[];
  failed_geocoded_orders?: FailedGeocodedOrder[];
  geocoded_orders?: GeocodedOrder[];
  error?: string;
};

export type OrderDelivery = GeocodedOrder &
  (
    | {
        has_associated_delivery: true;
        delivery_id: number;
        delivery_label: string;
        delivery_color: string;
      }
    | { has_associated_delivery: false }
  );
export type GeocodeOrdersRes = {
  ok: boolean;
  error?: string;
};

export type GetSortedOrdersRes = {
  ok: boolean;
  error?: string;
  orders: OrderDelivery[];
  deliveries: Delivery[];
};

export function getOrders() {
  return apiV1client.get<GetOrdersRes>("/api/v1/orders");
}

export type GetOrderRes = {
  ok: boolean;
  error?: string;
  order: GeocodedWithSuggestionsOrder;
};
export function getGeocodedWithSuggestionsOrder(id: string) {
  return apiV1client.get<GetOrderRes>("/api/v1/orders/single", {
    params: {
      order_id: id,
    },
  });
}
export function getSortedOrders() {
  return apiV1client.get<GetSortedOrdersRes>("/api/v1/orders/sorted");
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

export type ManualLocateRes = {
  ok: boolean;
  error?: string;
};
export function manuallyLocateOrder(data: FormData) {
  return apiV1client.patch<ManualLocateRes>("/api/v1/orders/geolocate", data);
}
export function manuallyLocateWithAddressUpdateOrder(data: FormData) {
  return apiV1client.patch<ManualLocateRes>(
    "/api/v1/orders/geolocate/update_address",
    data,
  );
}

export type UpdateOrderRes = {
  ok: boolean;
  error?: string;
}
export function updateOrder(data: FormData) {
  return apiV1client.patch<ManualLocateRes>("/api/v1/orders", data);
}
