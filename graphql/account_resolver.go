package main

import (
	"context"
	"log"
	"time"
)

type accountResolver struct {
	server *Server
}

func (a *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	ordersList, err := a.server.orderClient.GetOrdersByAccount(ctx, obj.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*Order

	for _, o := range ordersList {

		var products []*OrderProduct
		for _, p := range o.Products {

			products = append(products, &OrderProduct{
				ID:          p.Id,
				Description: p.Description,
				Name:        p.Name,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID:         o.Id,
			TotalPrice: o.TotalPrice,
			Products:   products,
		})
	}
	return orders, nil
}
