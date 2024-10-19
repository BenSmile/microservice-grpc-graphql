package account

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateAccount(context.Context, string) (*Account, error)
	GetAccountById(context.Context, string) (*Account, error)
	GetAccounts(context.Context, uint64, uint64) ([]Account, error)
}

type Account struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &accountService{r}
}

func (s *accountService) CreateAccount(c context.Context, name string) (*Account, error) {
	a := Account{
		Name: name,
		Id:   uuid.New().String(),
	}
	return s.repository.CreateAccount(c, a)
}

func (s *accountService) GetAccountById(c context.Context, id string) (*Account, error) {
	return s.repository.GetAccountById(c, id)
}

func (s *accountService) GetAccounts(c context.Context, skip, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccounts(c, skip, take)
}
