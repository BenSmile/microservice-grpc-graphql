package catalog

import (
	"context"
	"log"

	"github.com/bensmile/microservice-grpc-graphql/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	c := pb.NewCatalogServiceClient(conn)

	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateProduct(ctx context.Context, name, desc string, price float64) (*Product, error) {

	r, err := c.service.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:        name,
		Description: desc,
		Price:       price,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Product{
		Id:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, err
}

func (c *Client) GetProductById(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProductById(ctx, &pb.GetProductByIdRequest{Id: id})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Product{
		Id:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, err

}

func (c *Client) GetProducts(ctx context.Context, query string, ids []string, skip, take uint64) ([]Product, error) {
	r, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip:  skip,
		Take:  take,
		Ids:   ids,
		Query: query,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []Product{}

	for _, p := range r.Products {
		products = append(products, Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}
