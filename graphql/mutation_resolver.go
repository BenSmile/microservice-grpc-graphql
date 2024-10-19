package main

import "context"

type mutationResolver struct {
	server *Server
}

func (s *mutationResolver) CreateAccount(ctx context.Context, account *AccountInput) (*Account, error) {
	return nil, nil
}
func (s *mutationResolver) CreateProduct(ctx context.Context, product *ProductInput) (*Product, error) {
	return nil, nil
}
func (s *mutationResolver) CreateOrder(ctx context.Context, order *OrderInput) (*Order, error) {
	return nil, nil
}
