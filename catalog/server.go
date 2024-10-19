package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/bensmile/microservice-grpc-graphql/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

func ListenGrpc(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(serv, &grpcServer{pb.UnimplementedCatalogServiceServer{}, s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) CreateProduct(ctx context.Context, r *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	p, err := s.service.CreateProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}}, err
}
func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res []Product
	var err error
	if r.Query != "" {
		res, err = s.service.SearchProduct(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductsByIds(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []*pb.Product{}
	for _, p := range res {
		products = append(
			products,
			&pb.Product{
				Id:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
	}
	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}
func (s *grpcServer) GetProductById(ctx context.Context, r *pb.GetProductByIdRequest) (*pb.GetProductByIdResponse, error) {
	p, err := s.service.GetProductById(ctx, r.Id)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.GetProductByIdResponse{
		Product: &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, err
}
