package order

import (
	"context"
	"main/model"
)

type OrderRepository interface {
	Insert(ctx context.Context, order model.Order) error
	FindByID(ctx context.Context, id uint64) (model.Order, error)
	FindAll(ctx context.Context, page FindAllPage) (FindResult, error)
	Update(ctx context.Context, order model.Order) error
	DeleteByID(ctx context.Context, id uint64) error
}
