package geocode

import (
	"brotchaschte2/internal/database"
	query "brotchaschte2/internal/database/generated"
	"context"
	"database/sql"
	"log"
	"sync"
)

type GeocodeStore struct {
	db      *sql.DB
	queries *query.Queries
	mu      sync.Mutex
}

func New() GeocodeStore {
	db := database.New()
	return GeocodeStore{
		db:      db.DB(),
		queries: db.Queries(),
	}
}

// CheckActive checks whether currently a geocoding process is in progress
func (g *GeocodeStore) CheckActive() bool {
	// try to acquire lock. if successfull no geocoding is in progress
	success := g.mu.TryLock()
	if success {
		g.mu.Unlock()
		return false
	}
	return true
}

func (g *GeocodeStore) GeocodeOrders(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("failed to obtain db transaction for geocoding")
		return err
	}
	defer tx.Rollback()
	q := g.queries.WithTx(tx)

	orders, err := q.GetNotGeocodedOrdersIncludeFailed(ctx)
	// TODO: first try local geocoding

	addresses := make([]Address, len(orders))
	for i, order := range orders {
		addresses[i] = Address{
			id:      order.ID,
			address: order.Address,
			zipcode: order.Postalcode,
			city:    order.City,
		}
	}
	geocodes := AddressesToCoords(addresses, LV95)

	for _, geocode := range geocodes {
		if geocode.ok {
			q.CreateGeocode(ctx, query.CreateGeocodeParams{
				OrderID: geocode.Id,
				X:       geocode.X,
				Y:       geocode.Y,
				Label:   geocode.Label,
			})
			continue
		}

		if geocode.hasSuggestion {
			q.CreateGeocodeErrorWithSuggestion(ctx, query.CreateGeocodeErrorWithSuggestionParams{
				OrderID:           geocode.Id,
				Message:           geocode.error_message,
				SuggestedEasting:  sql.NullFloat64{Float64: geocode.X, Valid: true},
				SuggestedNorthing: sql.NullFloat64{Float64: geocode.Y, Valid: true},
				SuggestedLabel:    sql.NullString{String: geocode.Label, Valid: true},
			})
			continue
		}
		q.CreateGeocodeError(ctx, query.CreateGeocodeErrorParams{
			OrderID: geocode.Id,
			Message: geocode.error_message,
		})
	}
	if err := q.CleanUpErrors(ctx); err != nil {
		return err
	}

	return tx.Commit()

}

func (g *GeocodeStore) GetFailedGeocodeAttempts(ctx context.Context) ([]query.GeocodeError, error) {
	return g.queries.GetFailedGeocodes(ctx)
}
func (g *GeocodeStore) GetFailedGeocodeOrders(ctx context.Context) ([]query.GetFailedGeocodeOrdersRow, error) {
	return g.queries.GetFailedGeocodeOrders(ctx)
}
func (g *GeocodeStore) GetNotGeocodedIncludeFailed(ctx context.Context) ([]query.GetNotGeocodedOrdersIncludeFailedRow, error) {
	return g.queries.GetNotGeocodedOrdersIncludeFailed(ctx)
}
func (g *GeocodeStore) GetNotGeocodedExcludeFailed(ctx context.Context) ([]query.RawOrder, error) {
	return g.queries.GetNotGeocodedOrdersExcludeFailed(ctx)
}
func (g *GeocodeStore) GetGeocodedOrders(ctx context.Context) ([]query.GetGeocodedOrdersRow, error) {
	return g.queries.GetGeocodedOrders(ctx)
}

func (g *GeocodeStore) UpdateAddress(ctx context.Context, order_id, address, zip, city string) (query.RawOrder, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("failed to obtain db transaction for updating address")
		return query.RawOrder{}, err
	}
	defer tx.Rollback()
	q := g.queries.WithTx(tx)

	order, err := q.UpdateAddress(ctx, query.UpdateAddressParams{ID: order_id, Address: address, Postalcode: zip, City: city})
	if err != nil {
		return query.RawOrder{}, err
	}
	if err := q.DeleteGeocodeError(ctx, order_id); err != nil {
		return query.RawOrder{}, err
	}

	return order, tx.Commit()
}

func (g *GeocodeStore) ManuallyLocateOrder(ctx context.Context, order_id string, northing, easting float64) (query.OrderCoordinate, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("failed to obtain db transaction for geocoding")
		return query.OrderCoordinate{}, err
	}
	defer tx.Rollback()
	q := g.queries.WithTx(tx)

	ordercoord, err := q.ManualGeocode(ctx, query.ManualGeocodeParams{OrderID: order_id, X: northing, Y: easting})
	if err != nil {
		return query.OrderCoordinate{}, err
	}
	err = q.CleanUpErrors(ctx)
	if err != nil {
		return query.OrderCoordinate{}, err
	}

	return ordercoord, tx.Commit()
}

func (g *GeocodeStore) ManuallyLocateAndUpdateAddressOrder(ctx context.Context, order_id, address, postalcode, city string, northing, easting float64) (query.OrderCoordinate, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("failed to obtain db transaction for geocoding")
		return query.OrderCoordinate{}, err
	}
	defer tx.Rollback()
	q := g.queries.WithTx(tx)

	_, err = q.UpdateAddress(ctx, query.UpdateAddressParams{ID: order_id, Address: address, Postalcode: postalcode, City: city})
	if err != nil {
		return query.OrderCoordinate{}, err
	}

	ordercoord, err := q.ManualGeocode(ctx, query.ManualGeocodeParams{OrderID: order_id, X: northing, Y: easting})
	if err != nil {
		return query.OrderCoordinate{}, err
	}
	err = q.CleanUpErrors(ctx)
	if err != nil {
		return query.OrderCoordinate{}, err
	}

	return ordercoord, tx.Commit()
}

func (g *GeocodeStore) GetOrder(ctx context.Context, order_id string) (query.GetGeocodeWithSuggestionOrderRow, error) {
	return g.queries.GetGeocodeWithSuggestionOrder(ctx, order_id)
}
