import { createContext } from "svelte";
import type {
  GeocodedOrder,
  GeocodedWithSuggestionsOrder,
  OrderDelivery,
} from "@/lib/api/orders";
import { STANDARD_COLORS } from "./marker";
import { SvelteMap } from "svelte/reactivity";
import type { Delivery } from "@/lib/api/deliveries";

type DeliveryId = number;
type OrderId = string;
type Color = string;

type MarkerViewMode = "group-overview" | "delivery" | "manual";

export type MarkerEntry = {
  order: GeocodedOrder;
  delivery_id?: DeliveryId;
  color: Color;
  visible: boolean;
  rotation?: number;
};

type RestoreLogElement =
  | {
      type: "assignment";
      order_id: OrderId;
      old_delivery_id?: DeliveryId;
      new_delivery_id?: DeliveryId;
    }
  | {
      type: "add";
      entry: MarkerEntry;
    };

export class MarkerStore {
  // reactive map
  private markers: SvelteMap<OrderId, MarkerEntry>;
  private delivery_colors: Map<DeliveryId, Color>;
  private delivery_store: Map<DeliveryId | "unassigned", Set<OrderId>>;

  private viewMode: MarkerViewMode;
  private selected_delivery: DeliveryId | null;
  private restore_log: RestoreLogElement[];

  private callbacks: { marker_onclick: (order_id: OrderId) => void };

  public manual_order: GeocodedWithSuggestionsOrder | null;

  constructor() {
    this.markers = new SvelteMap();
    this.delivery_colors = new Map();
    this.delivery_store = new Map();
    this.viewMode = "group-overview";
    this.selected_delivery = null;
    this.manual_order = $state(null);
    this.restore_log = [];
    this.callbacks = { marker_onclick: () => {} };
  }

  public commitChanges() {
    this.restore_log = [];
  }
  public restoreChange() {
    const change = this.restore_log.pop();
    if (!change) {
      console.warn("nothing to restore");
      return;
    }
    if (change.type === "assignment") {
      if (!change.new_delivery_id) {
        if (!change.old_delivery_id) {
          // order was unassigned and got unassigned?
          console.warn(
            "not restoring: order was unassigned and got unassigned",
          );
          return;
        }
        // order got unassigned
        this.addOrderToDelivery(change.order_id, change.old_delivery_id);
        // ignore last change
        this.restore_log.pop();
        return;
      }
      // order was assigned
      this.unassignOrder(change.order_id);
      // ignore last change
      this.restore_log.pop();
      if (change.old_delivery_id) {
        // order was not unassigned previously
        this.addOrderToDelivery(change.order_id, change.old_delivery_id);
        // ignore last change
        this.restore_log.pop();
      }
    } else if (change.type === "add") {
      // order was temporarily added, thus needs to be removed
      const order_id = change.entry.order.id;
      const delivery_id = this.getOrderAssignment(change.entry.order.id);
      if (delivery_id) {
        this.unassignOrder(order_id);
        // ignore last change
        this.restore_log.pop();
      }
      this.markers.delete(order_id);
    }
  }
  public restoreChanges() {
    while (this.restore_log.length > 0) {
      this.restoreChange();
    }
  }

  get markerlist() {
    return this.markers;
  }

  public resetCallbacks() {
    this.callbacks = { marker_onclick: () => {} };
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
      this.delivery_store.set(delivery.id, new Set());
    }

    for (const order of orders) {
      if (order.has_associated_delivery) {
        this.delivery_store.get(order.delivery_id)?.add(order.id);
      } else {
        this.delivery_store.get("unassigned")?.add(order.id);
      }

      // i hate this ternary
      const delivery_id = order.has_associated_delivery
        ? order.delivery_id
        : undefined;

      const color = this.computeColor(delivery_id);
      const visible = this.viewMode !== "manual";
      this.markers.set(order.id, {
        color,
        delivery_id: delivery_id,
        order,
        visible,
      });
    }

    const overlaps = this.checkOverlaps();
    console.log(`found ${overlaps.length} overlaps`);
    for (const overlap of overlaps) {
      const rotation_unit = 33;
      let rotation = (-(overlap.size - 1) / 2) * rotation_unit;

      for (const id of overlap) {
        const m = this.markers.get(id)!;
        m.rotation = rotation;
        rotation += rotation_unit;
        this.markers.set(id, { ...m });
      }
    }
  }

  private computeColor(delivery_id?: DeliveryId): Color {
    switch (this.viewMode) {
      case "group-overview":
        if (!delivery_id) {
          return STANDARD_COLORS.unassigned;
        }
        const color = this.delivery_colors.get(delivery_id);
        if (!color) {
          console.warn(`delivery id: ${delivery_id} not in delivery_colormap`);
          return STANDARD_COLORS.unassigned;
        }
        return color;

      case "delivery":
        if (!this.selected_delivery) {
          console.warn(
            "MarkerViewMode set to 'delivery' but no delivery selected",
          );
        }
        if (!delivery_id) {
          return STANDARD_COLORS.unassigned;
        }
        if (delivery_id === this.selected_delivery) {
          const color = this.delivery_colors.get(delivery_id);
          if (!color) {
            console.warn(
              `delivery id: ${delivery_id} not in delivery_colormap`,
            );
            return STANDARD_COLORS.assigned;
          }
          return color;
        }
        return STANDARD_COLORS.unavailable;
      case "manual":
        return STANDARD_COLORS.unassigned;
    }
  }

  public setToDeliveryEditMode(selected_delivery: DeliveryId) {
    console.log(`change to delivery edit mode for ${selected_delivery}`);
    this.selected_delivery = selected_delivery;
    this.manual_order = null;
    this.viewMode = "delivery";

    for (const [delivery_id, order_ids] of this.delivery_store.entries()) {
      const color = this.computeColor(
        delivery_id === "unassigned" ? undefined : delivery_id,
      );
      for (const order_id of order_ids) {
        const marker = this.markers.get(order_id);
        if (!marker) {
          console.warn(
            `Order Id: '${order_id}' in delivery_store but not in markers`,
          );
          continue;
        }
        this.markers.set(order_id, { ...marker, color, visible: true });
      }
    }
  }

  public setToDeliveryOverviewMode() {
    console.log(`change to delivery overview mode`);
    this.selected_delivery = null;
    this.manual_order = null;
    this.viewMode = "group-overview";

    for (const [delivery_id, order_ids] of this.delivery_store.entries()) {
      const color = this.computeColor(
        delivery_id === "unassigned" ? undefined : delivery_id,
      );
      for (const order_id of order_ids) {
        const marker = this.markers.get(order_id);
        if (!marker) {
          console.warn(
            `Order Id: '${order_id}' in delivery_store but not in markers`,
          );
          continue;
        }
        this.markers.set(order_id, { ...marker, color, visible: true });
      }
    }
  }

  public setToManualLocateMode() {
    this.viewMode = "manual";

    for (const marker of this.markers.values()) {
      this.markers.set(marker.order.id, {
        ...marker,
        visible: false,
      });
    }
  }
  public addManualLocateOrder(order: GeocodedWithSuggestionsOrder) {
    this.manual_order = order;
  }

  public setDeliveryColor(delivery_id: DeliveryId, color: Color) {
    this.delivery_colors.set(delivery_id, color);
    const order_ids = this.delivery_store.get(delivery_id);
    if (!order_ids) {
      console.warn(
        `set delivery color for delivery_id: '${delivery_id}' but no orders in delivery_store`,
      );
      return;
    }
    for (const order_id of order_ids) {
      const marker = this.markers.get(order_id);
      if (!marker) {
        console.warn(
          `tried setting color for marker ${order_id}, but didn't exist in markers`,
        );
        continue;
      }
      this.markers.set(order_id, { ...marker, color });
    }
  }

  public addOrderToDelivery(order_id: OrderId, delivery_id: DeliveryId) {
    const ok = this.delivery_store.get("unassigned")?.delete(order_id);
    if (!ok) {
      console.warn(
        `adding order ${order_id} to delivery ${delivery_id} but order was not unassigned.`,
      );
    }
    this.delivery_store.get(delivery_id)?.add(order_id);
    const marker = this.markers.get(order_id);
    if (!marker) {
      console.warn(`order ${order_id} does not exist, cannot add to delivery`);
      return;
    }
    const color = this.computeColor(delivery_id);
    const new_marker = { ...marker, delivery_id, color };
    this.markers.set(order_id, new_marker);

    // log changes
    this.restore_log.push({
      type: "assignment",
      order_id,
      old_delivery_id: marker.delivery_id,
      new_delivery_id: delivery_id,
    });
    return new_marker;
  }

  public unassignOrder(order_id: OrderId) {
    const order = this.markers.get(order_id);
    if (!order) {
      console.warn(`while unassigning order ${order_id}: order not found`);
      return;
    }
    const current_delivery = order.delivery_id;
    if (!current_delivery) {
      console.warn("order already unassigned");
      return;
    }
    const orders = this.delivery_store.get(current_delivery);
    if (!orders) {
      console.warn(
        `unassign order ${order_id}: delivery ${current_delivery} not found`,
      );
      console.log(this.delivery_store);
      return;
    }
    orders.delete(order_id);
    order.delivery_id = undefined;
    this.delivery_store.get("unassigned")?.add(order_id);
    this.markers.set(order_id, {
      ...order,
      delivery_id: undefined,
      color: this.computeColor(undefined),
    });

    // log changes
    this.restore_log.push({
      type: "assignment",
      order_id,
      old_delivery_id: current_delivery,
    });
  }

  public getOrderAssignment(order_id: OrderId): DeliveryId | undefined {
    const order = this.markers.get(order_id);
    if (!order) {
      console.warn(`order ${order_id}: does not exist`);
      return;
    }
    return order.delivery_id;
  }

  public getOrdersForDelivery(delivery_id: DeliveryId) {
    const order_ids = this.delivery_store.get(delivery_id);
    if (!order_ids) {
      console.warn(`delivery ${delivery_id} doesn't exist`);
      return [];
    }
    const orders: GeocodedOrder[] = [];
    for (const order_id of order_ids) {
      const marker = this.markers.get(order_id);
      if (!marker) {
        console.warn(
          `while gathering orders for delivery ${delivery_id} order ${order_id} not found`,
        );
        continue;
      }
      orders.push(marker.order);
    }
    return orders;
  }

  public filterForDeliveryTime(time: "early" | "late" | "both") {
    for (const order of this.markers.values()) {
      let visible: boolean;
      switch (time) {
        case "early":
          visible = order.order.deliverymethod === "8 Uhr";
          break;
        case "late":
          visible = order.order.deliverymethod === "10 Uhr";
          break;
        case "both":
          visible = true;
          break;
      }
      this.markers.set(order.order.id, { ...order, visible });
    }
  }

  private checkOverlaps(): Set<OrderId>[] {
    // rather have some false positives!
    const DISTANCE = 8;
    const adj_map = new Map<OrderId, Set<OrderId>>();
    // compute distance
    // n squared babyy
    for (const [id1, m1] of this.markers.entries()) {
      for (const [id2, m2] of this.markers.entries()) {
        if (id1 === id2) {
          continue;
        }

        const distance =
          (m1.order.x - m2.order.x) ** 2 + (m1.order.y - m2.order.y) ** 2;

        if (distance <= DISTANCE) {
          if (!adj_map.has(id1)) {
            adj_map.set(id1, new Set());
          }
          if (!adj_map.has(id2)) {
            adj_map.set(id2, new Set());
          }
          adj_map.get(id1)!.add(id1);
          adj_map.get(id1)!.add(id2);
          adj_map.get(id2)!.add(id1);
          adj_map.get(id2)!.add(id2);
        }
      }
    }

    const overlaps: Set<OrderId>[] = [];
    const orders = Array.from(adj_map.keys());
    for (const order of orders) {
      if (!adj_map.has(order)) continue;
      overlaps.push(adj_map.get(order)!);
      adj_map.delete(order);
    }
    return overlaps;
  }
}

const [_getMarkerContext, _setMarkerContext] = createContext<MarkerStore>();

export function setMarkerStoreContext() {
  return _setMarkerContext(new MarkerStore());
}

export const getMarkerStoreContext = _getMarkerContext;
