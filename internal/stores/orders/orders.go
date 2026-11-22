package orders

import (
	"brotchaschte2/internal/database"
	query "brotchaschte2/internal/database/generated"
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type OrderStore struct {
	db      *sql.DB
	queries *query.Queries
}

func New() OrderStore {
	db := database.New()
	return OrderStore{
		db:      db.DB(),
		queries: db.Queries(),
	}
}

func (o *OrderStore) newCreateOrderParams(order query.RawOrder) query.CreateOrderParams {
	return query.CreateOrderParams{
		ID:             order.ID,
		Lastname:       order.Lastname,
		Firstname:      order.Firstname,
		Phonenumber:    order.Phonenumber,
		Email:          order.Email,
		Amountsizes:    order.Amountsizes,
		Amountsizem:    order.Amountsizem,
		Amountsizel:    order.Amountsizel,
		Address:        order.Address,
		City:           order.City,
		Postalcode:     order.Postalcode,
		Deliverymethod: order.Deliverymethod,
		Remarks:        order.Remarks,
		Ordermethod:    order.Ordermethod,
	}
}

func (o *OrderStore) GetOrders(ctx context.Context) ([]query.GetOrdersRow, error) {
	return o.queries.GetOrders(ctx)
}

func (o *OrderStore) GetStats(ctx context.Context) query.GetOrderStatsRow {
	stats, err := o.queries.GetOrderStats(ctx)
	if err != nil {
		return query.GetOrderStatsRow{}
	}
	return stats
}

func (o *OrderStore) Create(ctx context.Context, order query.RawOrder) error {
	return o.queries.CreateOrder(ctx, o.newCreateOrderParams(order))

}

func (o *OrderStore) CreateBulk(ctx context.Context, orders []query.RawOrder) error {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := o.queries.WithTx(tx)
	for _, order := range orders {
		q.CreateOrder(ctx, o.newCreateOrderParams(order))
	}
	return tx.Commit()
}

func (o *OrderStore) GenerateRandomOrders(ctx context.Context, n uint8) []query.RawOrder {
	addrs, err := o.queries.GetRandomAddresses(ctx, int64(n))
	if err != nil {
		panic(err)
	}

	orders := make([]query.RawOrder, n)

	for i, a := range addrs {

		amounts := [3]int64{rand.Int64N(2) + 1, 0, 0}
		rand.Shuffle(3, func(i, j int) {
			tmp := amounts[i]
			amounts[i] = amounts[j]
			amounts[j] = tmp
		})

		ziplabel := strings.Split(a.ZipLabel, " ")

		deliveryMethod := "8 Uhr"
		if rand.IntN(2) == 1 {
			deliveryMethod = "10 Uhr"
		}

		orders[i] = query.RawOrder{
			ID:             fmt.Sprintf("%d+%d", time.Now().UnixNano(), i),
			Lastname:       lastNames[rand.IntN(len(lastNames))],
			Firstname:      firstNames[rand.IntN(len(firstNames))],
			Phonenumber:    "-",
			Email:          "-",
			Amountsizes:    amounts[0],
			Amountsizem:    amounts[1],
			Amountsizel:    amounts[2],
			Address:        fmt.Sprintf("%s %s", a.StnLabel, a.AdrNumber),
			Postalcode:     ziplabel[0],
			City:           ziplabel[1],
			Deliverymethod: deliveryMethod,
			Remarks:        sql.NullString{},
			Ordermethod:    sql.NullString{},
			Vendorgroup:    sql.NullString{},
		}
	}

	return orders
}

func (o *OrderStore) SyncWithUpstream(ctx context.Context) error {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT \"id\", \"lastName\", \"firstName\", \"phoneNumber\", \"email\", \"amountSizeS\", \"amountSizeM\", \"amountSizeL\", \"address\", \"city\", \"postalCode\", \"deliveryMethod\", \"remarks\", \"orderMethod\", \"vendorGroup\" FROM \"Order\";")
	if err != nil {
		return fmt.Errorf("sync query failed: %w", err)
	}
	defer rows.Close()

	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("couldnt get tx: %w", err)
	}
	defer tx.Rollback()
	qtx := o.queries.WithTx(tx)

	for rows.Next() {
		var order query.CreateOrderIfNotExistParams

		err := rows.Scan(&order.ID, &order.Lastname, &order.Firstname, &order.Phonenumber, &order.Email, &order.Amountsizes, &order.Amountsizem, &order.Amountsizel, &order.Address, &order.City, &order.Postalcode, &order.Deliverymethod, &order.Remarks, &order.Ordermethod, &order.Vendorgroup)
		if err != nil {
			return fmt.Errorf("failed scan order: %w", err)
		}

		if err := qtx.CreateOrderIfNotExist(ctx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
	}

	return tx.Commit()
}

func (o *OrderStore) UpdateOrder(ctx context.Context,
	id string,
	lastname,
	firstname,
	phonenumber,
	email string,
	amountsizes,
	amountsizem,
	amountsizel int64,
	address,
	city,
	postalcode,
	deliverymethod string,
	remarks sql.NullString,
) error {
	return o.queries.UpdateOrder(ctx, query.UpdateOrderParams{
		ID:             id,
		Lastname:       lastname,
		Firstname:      firstname,
		Phonenumber:    phonenumber,
		Email:          email,
		Amountsizes:    amountsizes,
		Amountsizem:    amountsizem,
		Amountsizel:    amountsizel,
		Address:        address,
		City:           city,
		Postalcode:     postalcode,
		Deliverymethod: deliverymethod,
		Remarks:        remarks,
	})
}
