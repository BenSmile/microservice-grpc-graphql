package catalog

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type Service interface {
	CreateProduct(context.Context, string, string, float64) (*Product, error)
	GetProductById(context.Context, string) (*Product, error)
	GetProducts(context.Context, uint64, uint64) ([]Product, error)
	GetProductsByIds(context.Context, []string) ([]Product, error)
	SearchProduct(context.Context, string, uint64, uint64) ([]Product, error)
}

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{repository: r}
}

func (s *catalogService) CreateProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	p := Product{
		Id:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}
	if err := s.repository.CreateProduct(ctx, p); err != nil {
		log.Println(err)
		return nil, err
	}
	return &p, nil
}

func (s *catalogService) GetProductById(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductById(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.GetProducts(ctx, skip, take)
}

func (s *catalogService) GetProductsByIds(ctx context.Context, ids []string) ([]Product, error) {
	return s.repository.GetProductsByIds(ctx, ids)
}

func (s *catalogService) SearchProduct(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProduct(ctx, query, skip, take)
}
