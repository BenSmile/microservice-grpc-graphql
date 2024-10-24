package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/bensmile/microservice-grpc-graphql/order"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *Server
}

func (s *mutationResolver) CreateAccount(ctx context.Context, in *AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	newAccount, err := s.server.accountClient.CreateAccount(ctx, in.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Account{
		ID:   newAccount.Id,
		Name: newAccount.Name,
	}, nil
}

func (s *mutationResolver) CreateProduct(ctx context.Context, in *ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	newProduct, err := s.server.catalogClient.CreateProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Product{
		ID:          newProduct.Id,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
	}, nil
}

func (s *mutationResolver) CreateOrder(ctx context.Context, in *OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var orderedProducts []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		orderedProducts = append(orderedProducts, order.OrderedProduct{
			Id:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}
	newOrder, err := s.server.orderClient.PlaceOrder(ctx, in.AccountID, orderedProducts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []*OrderProduct
	for _, p := range newOrder.Products {
		products = append(products, &OrderProduct{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    int(p.Quantity),
		})
	}
	return &Order{
		ID:         newOrder.Id,
		CreatedAt:  newOrder.CreatedAt,
		TotalPrice: newOrder.TotalPrice,
		Products:   products,
	}, nil
}
