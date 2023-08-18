package usecase

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/vemta/api/domain/repository"
	uow "github.com/vemta/api/pkg"
	"github.com/vemta/common/entity"
)

type OrderCreateUsecaseInput struct {
	Customer      string   `json:"customer"`
	Items         []string `json:"items"`
	PaymentMethod string   `json:"payment_method"`
	DiscountCode  string   `json:"discount_code"`
}

type OrderCreateUsecase struct {
	Uow uow.UowInterface
}

func NewOrderCreateUsecase(uow uow.UowInterface) *OrderCreateUsecase {
	return &OrderCreateUsecase{
		Uow: uow,
	}
}

func (u *OrderCreateUsecase) Execute(ctx context.Context, input OrderCreateUsecaseInput) (*entity.Plugin, error) {
	orderRepository := u.getOrderRepository(ctx)

	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	id := strings.Replace(uid.String(), "-", "", 0)

	items := make([]entity.Item, 0, len(input.Items))
	for _, item := range input.Items {
		items = append(items, entity.Item{ID: item})
	}

	plugin := &entity.Order{
		ID: id,
		Customer: &entity.User{
			Email: input.Customer,
		},
		Items: &items,
	}

	message, err := json.Marshal(plugin)
	if err != nil {
		return nil, err
	}

	if err := producer.Produce("createPlugin", message); err != nil {
		return nil, err
	}

	return plugin, nil
}

func (u *OrderCreateUsecase) getOrderRepository(ctx context.Context) repository.OrderRepositoryInterface {
	orderRepository, err := u.Uow.GetRepository(ctx, "PluginRepository")
	if err != nil {
		panic(err)
	}
	return orderRepository.(repository.OrderRepositoryInterface)
}
