package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/bensmile/microservice-grpc-graphql/account"
)

type Server struct {
	accountClient *account.Client
	// catalogClient *catalog.Client
	// orderClient   *order.Client
}

func NewGraphqlServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {

	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}
	// catalogClient, err := catalog.NewClient(catalogUrl)
	// if err != nil {
	// 	catalogClient.Close()
	// 	return nil, err
	// }
	// orderClient, err := order.NewClient(orderUrl)
	// if err != nil {
	// 	orderClient.Close()
	// 	return nil, err
	// }

	return &Server{
		accountClient,
		// catalogClient,
		// orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
