package order

import (
	"context"
	"log"
	"time"

	"github.com/bensmile/microservice-grpc-graphql/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)

	return &Client{
		conn, c,
	}, nil
}

func (c *Client) PlaceOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error) {
	protoProducts := []*pb.PlaceOrderRequest_OrderedProduct{}

	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PlaceOrderRequest_OrderedProduct{
			ProductId: p.Id,
			Quantity:  p.Quantity,
		})
	}

	r, err := c.service.PlaceOrder(ctx,
		&pb.PlaceOrderRequest{
			AccountId: accountId,
			Products:  protoProducts,
		})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orderedProducts := []OrderedProduct{}

	for _, p := range r.Order.Products {
		orderedProducts = append(orderedProducts, OrderedProduct{
			Id:          p.Id,
			Quantity:    p.Quantity,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	createdAt := time.Time{}
	createdAt.UnmarshalBinary(r.Order.CreatedAt)

	return &Order{
		Id:         r.Order.Id,
		CreatedAt:  createdAt,
		TotalPrice: r.Order.TotalPrice,
		AccountId:  accountId,
		Products:   orderedProducts,
	}, nil
}

func (c *Client) GetOrdersByAccount(ctx context.Context, accountId string) ([]*Order, error) {

	r, err := c.service.GetOrdersAccountById(ctx,
		&pb.GetOrdersByAccountRequest{
			AccountId: accountId,
		})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orders := []*Order{}

	for _, o := range r.Orders {

		orderedProducts := []OrderedProduct{}

		for _, op := range o.Products {
			orderedProducts = append(orderedProducts, OrderedProduct{
				Id:          op.Id,
				Quantity:    op.Quantity,
				Name:        op.Name,
				Description: op.Description,
				Price:       op.Price,
			})
		}

		createdAt := time.Time{}
		createdAt.UnmarshalBinary(o.CreatedAt)

		orders = append(orders,
			&Order{
				Id:         o.Id,
				CreatedAt:  createdAt,
				TotalPrice: o.TotalPrice,
				AccountId:  accountId,
				Products:   orderedProducts,
			})
	}

	return orders, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
