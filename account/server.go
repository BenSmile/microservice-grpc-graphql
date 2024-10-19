package account

import (
	"context"
	"fmt"
	"net"

	"github.com/bensmile/microservice-grpc-graphql/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
}

func ListenGrpc(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) CreateAccount(c context.Context, r *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	a, err := s.service.CreateAccount(c, r.Name)
	if err != nil {
		return nil, err
	}

	return &pb.CreateAccountResponse{Account: &pb.Account{Id: a.Id, Name: a.Name}}, nil
}

func (s *grpcServer) GetAccountById(c context.Context, r *pb.GetAccountByIdRequest) (*pb.GetAccountByIdResponse, error) {
	a, err := s.service.GetAccountById(c, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountByIdResponse{Account: &pb.Account{Id: a.Id, Name: a.Name}}, nil
}

func (s *grpcServer) GetAccounts(c context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(c, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}

	for _, p := range res {
		accounts = append(accounts, &pb.Account{
			Id:   p.Id,
			Name: p.Name,
		})
	}

	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}
