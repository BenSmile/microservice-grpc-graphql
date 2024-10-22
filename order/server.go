package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/bensmile/microservice-grpc-graphql/account"
	"github.com/bensmile/microservice-grpc-graphql/catalog"
	"github.com/bensmile/microservice-grpc-graphql/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	*pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountUrl, catalogUrl string, port uint32) error {

	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		accountClient.Close()
		return err
	}

	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		catalogClient.Close()
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		&pb.UnimplementedOrderServiceServer{},
		s,
		accountClient,
		catalogClient,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PlaceOrder(ctx context.Context, r *pb.PlaceOrderRequest) (*pb.PlaceOrderReponse, error) {

	_, err := s.accountClient.GetAccountById(ctx, r.AccountId)

	if err != nil {
		log.Printf("Error when getting account : %v\n", err)
		return nil, errors.New("account not found")
	}

	productIds := []string{}
	orderedProducts := []OrderedProduct{}

	for _, p := range r.Products {
		orderedProducts = append(orderedProducts, OrderedProduct{
			Id:       p.ProductId,
			Quantity: p.Quantity,
		})
		productIds = append(productIds, p.ProductId)
	}

	products, err := s.catalogClient.GetProducts(ctx, "", productIds, 0, 0)
	if err != nil {
		log.Printf("Error when getting products : %v\n", err)
		return nil, errors.New("product not found")
	}
	if len(products) != len(r.Products) {
		return nil, errors.New("missing products")
	}

	for _, p := range products {
		op := OrderedProduct{
			Id:          p.Id,
			Price:       p.Price,
			Quantity:    0,
			Name:        p.Name,
			Description: p.Description,
		}
		for _, rp := range r.Products {
			if rp.ProductId == p.Id {
				op.Quantity = rp.Quantity
				break
			}
		}
		if op.Quantity != 0 {
			orderedProducts = append(orderedProducts, op)
		}
	}

	order, err := s.service.PlaceOrder(ctx, r.AccountId, orderedProducts)
	if err != nil {
		log.Println(err)
		return nil, errors.New("could not place order")
	}

	orderProto := &pb.Order{
		Id:         order.Id,
		AccountId:  order.AccountId,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderedProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products,
			&pb.Order_OrderedProduct{
				Id:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			},
		)
	}
	return &pb.PlaceOrderReponse{
		Order: orderProto,
	}, nil
}
func (s *grpcServer) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.GetOrderReponse, error) {
	return nil, nil
}

func (s *grpcServer) GetOrdersAccountById(ctx context.Context, p *pb.GetOrdersByAccountRequest) (*pb.GetOrdersByAccountReponse, error) {
	accountOrders, err := s.service.GetOrdersByAccount(ctx, p.AccountId)
	if err != nil {
		return nil, err
	}

	productsIdMap := map[string]bool{}

	for _, o := range accountOrders {
		for _, p := range o.Products {
			productsIdMap[p.Id] = true
		}
	}

	productIds := []string{}
	for id := range productsIdMap {
		productIds = append(productIds, id)
	}

	products, err := s.catalogClient.GetProducts(ctx, "", productIds, 0, 0)
	if err != nil {
		log.Printf("Error when getting products : %v\n", err)
		return nil, errors.New("product not found")
	}

	// *415*1#

	orders := []*pb.Order{}

	for _, o := range accountOrders {
		op := &pb.Order{
			AccountId:  o.AccountId,
			Id:         o.Id,
			TotalPrice: o.TotalPrice,
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()
		for _, product := range o.Products {
			for _, p := range products {
				if p.Id == product.Id {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}
			op.Products = append(op.Products, &pb.Order_OrderedProduct{
				Id:          product.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}
		orders = append(orders, op)
	}

	return &pb.GetOrdersByAccountReponse{
		Orders: orders,
	}, nil
}
