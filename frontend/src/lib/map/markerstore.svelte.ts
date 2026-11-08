import { createContext } from "svelte";
import type { GeocodedOrder, OrderDelivery } from "@/lib/api/orders";
import { STANDARD_COLORS } from "./marker";
import { SvelteMap } from "svelte/reactivity";

type DeliveryId = number;
type OrderId = string;
type Color = string;

type MarkerViewMode = "group-overview" | "delivery"

type MarkerEntry = {
	order: GeocodedOrder
	deliver_id?: DeliveryId
	color: Color
}

export class MarkerStore {
	// reactive map
	private markers: SvelteMap<OrderId, MarkerEntry>
	private delivery_colors: Map<DeliveryId, Color>
	private delivery_store: Map<DeliveryId | "unassigned", OrderId[]>

	private viewMode: MarkerViewMode
	private selected_delivery: DeliveryId | null

	constructor() {
		this.markers = new SvelteMap();
		this.delivery_colors = new Map();
		this.delivery_store = new Map();
		this.viewMode = "group-overview";
		this.selected_delivery = null;
	}

	get markerlist() {
		return this.markers;
	}

	public setOrders(orders: OrderDelivery[]) {
		this.markers.clear();
		this.delivery_colors.clear();
		this.delivery_store.clear();

		this.delivery_store.set("unassigned", new Array());

		for (const order of orders) {
			if (order.has_associated_delivery) {
				if (!this.delivery_colors.has(order.delivery_id)) {
					this.delivery_colors.set(order.delivery_id, order.delivery_color);
				}
				if (!this.delivery_store.has(order.delivery_id)) {
					this.delivery_store.set(order.delivery_id, new Array());
				}
				this.delivery_store.get(order.delivery_id)?.push(order.id);
			} else {
				this.delivery_store.get("unassigned")?.push(order.id);
			}

			// i hate this ternary
			const delivery_id = order.has_associated_delivery ? order.delivery_id : undefined;

			const color = this.computeColor(delivery_id);
			this.markers.set(order.id, {
				color,
				deliver_id: delivery_id,
				order
			});
		}
	}

	private computeColor(delivery_id?: DeliveryId): Color {
		switch (this.viewMode) {
			case "group-overview":
				if (!delivery_id) {
					return STANDARD_COLORS.unassigned
				}
				const color = this.delivery_colors.get(delivery_id)
				if (!color) {
					console.warn(`delivery id: ${delivery_id} not in delivery_colormap`);
					return STANDARD_COLORS.unassigned;
				}
				return color;

			case "delivery":
				if (!this.selected_delivery) {
					console.warn("MarkerViewMode set to 'delivery' but no delivery selected");
				}
				if (!delivery_id) {
					return STANDARD_COLORS.unassigned
				}
				if (delivery_id === this.selected_delivery) {
					return STANDARD_COLORS.assigned
				}
				return STANDARD_COLORS.unavailable
		}

	}

	public setToDeliveryEditMode(selected_delivery: DeliveryId) {
		console.log(`change to delivery edit mode for ${selected_delivery}`);
		this.selected_delivery = selected_delivery;
		this.viewMode = "delivery";

		for (const [delivery_id, order_ids] of this.delivery_store.entries()) {
			const color = this.computeColor(delivery_id === "unassigned" ? undefined : delivery_id);
			for (const order_id of order_ids) {
				const marker = this.markers.get(order_id);
				if (!marker) {
					console.warn(`Order Id: '${order_id}' in delivery_store but not in markers`);
					continue;
				}
				this.markers.set(order_id, { ...marker, color })
			}

		}
	}

	public setToDeliveryOverviewMode() {
		console.log(`change to delivery overview mode`);
		this.selected_delivery = null;
		this.viewMode = "group-overview";

		for (const [delivery_id, order_ids] of this.delivery_store.entries()) {
			const color = this.computeColor(delivery_id === "unassigned" ? undefined : delivery_id);
			for (const order_id of order_ids) {
				const marker = this.markers.get(order_id);
				if (!marker) {
					console.warn(`Order Id: '${order_id}' in delivery_store but not in markers`);
					continue;
				}
				marker.color = color;
				this.markers.set(order_id, marker);
			}

		}
	}

	public setDeliveryColor(delivery_id: DeliveryId, color: Color) {
		this.delivery_colors.set(delivery_id, color);
		const order_ids = this.delivery_store.get(delivery_id);
		if (!order_ids) {
			console.warn(`set delivery color for delivery_id: '${delivery_id}' but no orders in delivery_store`);
			return;
		}
		for (const order_id of order_ids) {
			const marker = this.markers.get(order_id);
			if (!marker) {
				console.warn(`tried setting color for marker ${order_id}, but didn't exist in markers`);
				continue;
			}
			marker.color = color;
			this.markers.set(order_id, marker);
		}

	}

}



const [_getMarkerContext, _setMarkerContext] = createContext<MarkerStore>();

export function setMarkerStoreContext() {
	return _setMarkerContext(new MarkerStore())
}

export const getMarkerStoreContext = _getMarkerContext
