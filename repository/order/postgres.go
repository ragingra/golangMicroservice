package order

import (
	"context"
	"database/sql"
	"fmt"
	"main/model"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	DB *sql.DB
}

var _ OrderRepository = (*PostgresRepo)(nil)

func (r *PostgresRepo) Insert(ctx context.Context, order model.Order) error {
	query := `INSERT INTO orders (order_id, customer_id, line_items, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.ExecContext(ctx, query, order.OrderID, order.CustomerID, order.LineItems, order.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}
	return nil
}

func (r *PostgresRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	var order model.Order
	query := `SELECT order_id, customer_id, line_items, created_at FROM orders WHERE order_id = $1`
	row := r.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&order.OrderID, &order.CustomerID, &order.LineItems, &order.CreatedAt)
	if err == sql.ErrNoRows {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("failed to find order: %w", err)
	}
	return order, nil
}

func (r *PostgresRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	query := `SELECT order_id, customer_id, line_items, created_at FROM orders OFFSET $1 LIMIT $2`
	rows, err := r.DB.QueryContext(ctx, query, page.Offset, page.Size)
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to find all orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderID, &order.CustomerID, &order.LineItems, &order.CreatedAt); err != nil {
			return FindResult{}, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return FindResult{}, fmt.Errorf("rows iteration failed: %w", err)
	}

	return FindResult{Orders: orders, Cursor: page.Offset + uint64(len(orders))}, nil
}

func (r *PostgresRepo) Update(ctx context.Context, order model.Order) error {
	query := `UPDATE orders SET customer_id = $2, line_items = $3, shipped_at = $4, completed_at = $5 WHERE order_id = $1`
	_, err := r.DB.ExecContext(ctx, query, order.OrderID, order.CustomerID, order.LineItems, order.ShippedAt, order.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

func (r *PostgresRepo) DeleteByID(ctx context.Context, id uint64) error {
	query := `DELETE FROM orders WHERE order_id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err == sql.ErrNoRows {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
