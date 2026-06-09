package delivery

import (
	"brotchaschte2/internal/database"
	query "brotchaschte2/internal/database/generated"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
)

type DeliveryStore struct {
	db      *sql.DB
	queries *query.Queries
}

func New() DeliveryStore {
	db := database.New()
	return DeliveryStore{
		db:      db.DB(),
		queries: db.Queries(),
	}
}

// GetOrderCoordinates  retrieves all orders with their coordinates.
// The result is sorted in ascending order of the northing coordinate. This prevents overlapping markers in the map
func (d *DeliveryStore) GetOrderCoordinates(ctx context.Context) ([]query.GetOrderCoordinatesRow, error) {
	return d.queries.GetOrderCoordinates(ctx)
}
func randomColorExceptYellow() string {
	return fmt.Sprintf("hsl(%d, 70%%, %d%%)", (rand.Intn(300)+60+30)%360, rand.Intn(21)+50)
}
func (d *DeliveryStore) CreateNewDelivery(ctx context.Context, label, until string) (query.Delivery, error) {
	return d.queries.CreateDelivery(ctx, query.CreateDeliveryParams{Label: label, Until: until, Color: randomColorExceptYellow()})
}

func (d *DeliveryStore) NewDeliveryColor(ctx context.Context, delivery_id int64) (query.Delivery, error) {
	return d.queries.UpdateDeliveryColor(ctx, query.UpdateDeliveryColorParams{ID: delivery_id, Color: randomColorExceptYellow()})
}

func (d *DeliveryStore) GetDelivery(ctx context.Context, id int64) (query.Delivery, error) {
	return d.queries.GetDelivery(ctx, id)
}

func (d *DeliveryStore) GetDeliveries(ctx context.Context) ([]query.Delivery, error) {
	return d.queries.GetDeliveries(ctx)
}

func (d *DeliveryStore) UpdateDelivery(ctx context.Context, id int64, until, label string) (query.Delivery, error) {
	return d.queries.UpdateDelivery(ctx, query.UpdateDeliveryParams{ID: id, Until: until, Label: label})
}

func (d *DeliveryStore) AssignOrders(ctx context.Context, id int64, selected_orders, notselected_orders []string) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queries := d.queries.WithTx(tx)
	defer tx.Rollback()

	for _, orderid := range notselected_orders {
		err := queries.DeleteDeliveryOrderDeselected(ctx, query.DeleteDeliveryOrderDeselectedParams{OrderID: orderid, DeliveryID: id})
		if err != nil {
			log.Printf("error while deleting deselected: %s\n", err.Error())
			continue
		}
	}

	for _, orderid := range selected_orders {
		err := queries.CreateDeliveryOrderSelected(ctx, query.CreateDeliveryOrderSelectedParams{OrderID: orderid, DeliveryID: id})
		if err != nil {
			log.Printf("error while updating selected: %s\n", err.Error())
			continue
		}
	}

	return tx.Commit()
}
