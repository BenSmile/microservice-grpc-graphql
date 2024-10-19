package account

import (
	"context"

	"github.com/bensmile/microservice-grpc-graphql/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	c := pb.NewAccountServiceClient(conn)

	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.CreateAccount(ctx,
		&pb.CreateAccountRequest{
			Name: name,
		})
	if err != nil {
		return nil, err
	}
	return &Account{
		Id:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccountById(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccountById(ctx,
		&pb.GetAccountByIdRequest{
			Id: id,
		})
	if err != nil {
		return nil, err
	}
	return &Account{
		Id:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{
		Skip: skip,
		Take: take,
	})
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			Id:   a.Id,
			Name: a.Name,
		})
	}
	return accounts, nil
}
