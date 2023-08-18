package repository

import (
	"context"

	"github.com/vemta/common/entity"
)

type OrderRepositoryInterface interface {
	Find(ctx context.Context, id string) (*entity.Order, error)
	FindUserOrder(ctx context.Context, userIdentifier string) (*[]entity.Order, error)
	Create(ctx context.Context, login *entity.Order) error
}
