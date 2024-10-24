package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (q *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil {
		r, err := q.server.accountClient.GetAccountById(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{{
			ID:   r.Id,
			Name: r.Name,
		}}, nil
	}

	skip, take := pagination.bound()

	accountsList, err := q.server.accountClient.GetAccounts(ctx, skip, take)

	if err != nil {
		return nil, err
	}

	var accounts []*Account

	for _, a := range accountsList {
		accounts = append(accounts,
			&Account{
				ID:   a.Id,
				Name: a.Name,
			})
	}

	return accounts, nil
}
func (q *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil {
		r, err := q.server.catalogClient.GetProductById(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{{
			ID:          r.Id,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		}}, nil
	}

	skip, take := pagination.bound()

	queryString := ""
	if query != nil {
		queryString = *query
	}
	productsList, err := q.server.catalogClient.GetProducts(ctx, queryString, nil, skip, take)
	if err != nil {
		return nil, err
	}
	var products []*Product

	for _, p := range productsList {
		products = append(products, &Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return nil, nil
}

func (p *PaginationInput) bound() (uint64, uint64) {
	skip, take := uint64(0), uint64(0)
	if p != nil {
		skip, take = uint64(*p.Skip), uint64(*p.Take)
	}
	return skip, take
}
