import { createContext } from "svelte";
import type { GeocodedOrder, OrderDelivery } from "@/lib/api/orders";
import { STANDARD_COLORS } from "./marker";
import { SvelteMap } from "svelte/reactivity";
import type { Delivery } from "@/lib/api/deliveries";

type DeliveryId = number;
type OrderId = string;
type Color = string;

type MarkerViewMode = "group-overview" | "delivery"

type MarkerEntry = {
	order: GeocodedOrder
	delivery_id?: DeliveryId
	color: Color
}

export class MarkerStore {
	// reactive map
	private markers: SvelteMap<OrderId, MarkerEntry>
	private delivery_colors: Map<DeliveryId, Color>
	private delivery_store: Map<DeliveryId | "unassigned", Set<OrderId>>

	private viewMode: MarkerViewMode
	private selected_delivery: DeliveryId | null
	private restore_log: { order_id: OrderId, old_delivery_id?: DeliveryId, new_delivery_id?: DeliveryId }[]

	private callbacks: { marker_onclick: (order_id: OrderId) => void }

	constructor() {
		this.markers = new SvelteMap();
		this.delivery_colors = new Map();
		this.delivery_store = new Map();
		this.viewMode = "group-overview";
		this.selected_delivery = null;
		this.restore_log = [];
		this.callbacks = { marker_onclick: () => { } }

	}

	public commitChanges() {
		this.restore_log = [];
	}
	public restoreChange() {
		const change = this.restore_log.pop();
		if (!change) {
			console.warn("nothing to restore")
			return;
		}
		if (!change.new_delivery_id) {
			if (!change.old_delivery_id) {
				// order was unassigned and got unassigned?
				console.warn("not restoring: order was unassigned and got unassigned")
				return
			}
			// order got unassigned
			this.addOrderToDelivery(change.order_id, change.old_delivery_id)
			// ignore last change
			this.restore_log.pop();
			return;
		}
		// order was assigned
		this.unassignOrder(change.order_id)
		// ignore last change
		this.restore_log.pop();
		if (change.old_delivery_id) {
			// order was not unassigned previously
			this.addOrderToDelivery(change.order_id, change.old_delivery_id)
			// ignore last change
			this.restore_log.pop();
		}
	}
	public restoreChanges() {
		while (this.restore_log.length > 0) {
			this.restoreChange()
		}
	}

	get markerlist() {
		return this.markers;
	}

	public resetCallbacks() {
		this.callbacks = { marker_onclick: () => { } }
	}

	public _markerOnClick(order_id: OrderId) {
		this.callbacks.marker_onclick(order_id);
	}
	public setMarkerOnClick(fn: (order_id: OrderId) => void) {
		this.callbacks.marker_onclick = fn;
	}

	public setOrders(orders: OrderDelivery[], deliveries: Delivery[]) {
		this.markers.clear();
		this.delivery_colors.clear();
		this.delivery_store.clear();

		this.delivery_store.set("unassigned", new Set());
		for (const delivery of deliveries) {
			this.delivery_colors.set(delivery.id, delivery.color);
			this.delivery_store.set(delivery.id, new Set())
		}

		for (const order of orders) {
			if (order.has_associated_delivery) {
				this.delivery_store.get(order.delivery_id)?.add(order.id);
			} else {
				this.delivery_store.get("unassigned")?.add(order.id);
			}

			// i hate this ternary
			const delivery_id = order.has_associated_delivery ? order.delivery_id : undefined;

			const color = this.computeColor(delivery_id);
			this.markers.set(order.id, {
				color,
				delivery_id: delivery_id,
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
					return STANDARD_COLORS.unassigned;
				}
				if (delivery_id === this.selected_delivery) {
					const color = this.delivery_colors.get(delivery_id)
					if (!color) {
						console.warn(`delivery id: ${delivery_id} not in delivery_colormap`);
						return STANDARD_COLORS.assigned;
					}
					return color;
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
				this.markers.set(order_id, { ...marker, color });
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
			this.markers.set(order_id, { ...marker, color });
		}

	}

	public addOrderToDelivery(order_id: OrderId, delivery_id: DeliveryId) {
		const ok = this.delivery_store.get("unassigned")?.delete(order_id);
		if (!ok) {
			console.warn(`adding order ${order_id} to delivery ${delivery_id} but order was not unassigned.`)
		}
		this.delivery_store.get(delivery_id)?.add(order_id);
		const marker = this.markers.get(order_id);
		if (!marker) {
			console.warn(`order ${order_id} does not exist, cannot add to delivery`);
			return;
		}
		const color = this.computeColor(delivery_id);
		const new_marker = { ...marker, delivery_id, color }
		this.markers.set(order_id, new_marker);

		// log changes
		this.restore_log.push({ order_id, old_delivery_id: marker.delivery_id, new_delivery_id: delivery_id })
		return new_marker;
	}

	public unassignOrder(order_id: OrderId) {
		const order = this.markers.get(order_id);
		if (!order) {
			console.warn(`while unassigning order ${order_id}: order not found`)
			return;
		}
		const current_delivery = order.delivery_id;
		if (!current_delivery) {
			console.warn("order already unassigned");
			return;
		}
		const orders = this.delivery_store.get(current_delivery)
		if (!orders) {
			console.warn(`unassign order ${order_id}: delivery ${current_delivery} not found`);
			console.log(this.delivery_store)
			return;
		}
		orders.delete(order_id);
		order.delivery_id = undefined;
		this.delivery_store.get("unassigned")?.add(order_id);
		this.markers.set(order_id, { ...order, delivery_id: undefined, color: this.computeColor(undefined) })

		// log changes
		this.restore_log.push({ order_id, old_delivery_id: current_delivery })
	}

	public getOrderAssignment(order_id: OrderId): DeliveryId | undefined {
		const order = this.markers.get(order_id);
		if (!order) {
			console.warn(`order ${order_id}: does not exist`)
			return;
		}
		return order.delivery_id
	}


	public getOrdersForDelivery(delivery_id: DeliveryId) {
		const order_ids = this.delivery_store.get(delivery_id);
		if (!order_ids) {
			console.warn(`delivery ${delivery_id} doesn't exist`)
			return [];
		}
		const orders: GeocodedOrder[] = [];
		for (const order_id of order_ids) {
			const marker = this.markers.get(order_id)
			if (!marker) {
				console.warn(`while gathering orders for delivery ${delivery_id} order ${order_id} not found`)
				continue;
			}
			orders.push(marker.order)
		}
		return orders;

	}

}



const [_getMarkerContext, _setMarkerContext] = createContext<MarkerStore>();

export function setMarkerStoreContext() {
	return _setMarkerContext(new MarkerStore())
}

export const getMarkerStoreContext = _getMarkerContext
