package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PlaceOrder(context.Context, Order) error
	GetOrdersByAccount(context.Context, string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {

}
func (r *postgresRepository) PlaceOrder(ctx context.Context, order Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders(id, created_at, account_id, total_price)
			VALUES ($1, $2, $3, $4)`,
		order.Id,
		order.CreatedAt,
		order.AccountId,
		order.TotalPrice)

	if err != nil {
		return err
	}

	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))

	for _, p := range order.Products {
		_, err = stmt.ExecContext(ctx, order.Id, p.ProductId, p.Quantity)
		if err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx)

	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}

func (r *postgresRepository) GetOrdersByAccount(ctx context.Context, accountId string) ([]Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT 
			o.id,
			o.created_at,
			o.account_id,
			o.total_price::money::numeric::float8 
			op.product_id,
			op.quantity
			FROM orders o
			JOIN order_products op ON (o.id=op.order_id)
			WHERE o.account_id = $1 
			ORDER BY o.id`,
		accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	var products []OrderedProduct

	for rows.Next() {
		var order Order
		if err := rows.Scan(
			&order.Id,
			&order.CreatedAt,
			&order.AccountId,
			&orderedProduct.ProductId,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}
		if lastOrder.Id != "" && lastOrder.Id != order.Id {
			newOrder := Order{
				Id:         lastOrder.Id,
				AccountId:  lastOrder.AccountId,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   lastOrder.Products,
			}
			orders = append(orders, newOrder)
			products = []OrderedProduct{}
		}
		products = append(products, OrderedProduct{
			ProductId: orderedProduct.ProductId,
			Quantity:  orderedProduct.Quantity,
		})
		*lastOrder = order
	}
	if lastOrder != nil {
		newOrder := Order{
			Id:         lastOrder.Id,
			AccountId:  lastOrder.AccountId,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   lastOrder.Products,
		}
		orders = append(orders, newOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
