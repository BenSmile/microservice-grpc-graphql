// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package main

import (
	"time"
)

type AccountInput struct {
	Name string `json:"name"`
}

type Mutation struct {
}

type Order struct {
	ID         string          `json:"id"`
	CreatedAt  time.Time       `json:"createdAt"`
	TotalPrice float64         `json:"totalPrice"`
	Products   []*OrderProduct `json:"products"`
}

type OrderInput struct {
	AccountID string               `json:"accountId"`
	Products  []*OrderProductInput `json:"products"`
}

type OrderProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type OrderProductInput struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type PaginationInput struct {
	Skip *int `json:"skip,omitempty"`
	Take *int `json:"take,omitempty"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Query struct {
}