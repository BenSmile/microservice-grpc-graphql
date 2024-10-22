package order

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	PlaceOrder(context.Context, string, []OrderedProduct) (*Order, error)
	GetOrdersByAccount(context.Context, string) ([]Order, error)
}

type Order struct {
	Id         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountId  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	Id          string
	Quantity    uint32
	Name        string
	Description string
	Price       float64
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s *orderService) PlaceOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error) {
	o := &Order{
		Id:        uuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountId: accountId,
		Products:  products,
	}

	o.TotalPrice = 0

	for _, p := range products {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}
	if err := s.repository.PlaceOrder(ctx, *o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *orderService) GetOrdersByAccount(ctx context.Context, accountId string) ([]Order, error) {
	return s.repository.GetOrdersByAccount(ctx, accountId)
}
