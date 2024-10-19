package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	Ping() error
	CreateAccount(context.Context, Account) (*Account, error)
	GetAccountById(context.Context, string) (*Account, error)
	ListAccounts(context.Context, uint64, uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) CreateAccount(ctx context.Context, a Account) (*Account, error) {
	if err := r.db.QueryRowContext(ctx, "INSERT INTO accounts(id, name) VALUES ($1,$2) RETURNING id, name",
		a.Id, a.Name).Scan(&a.Id, &a.Name); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *postgresRepository) GetAccountById(ctx context.Context, id string) (*Account, error) {
	var account Account
	if err := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id).
		Scan(&account.Id, &account.Name); err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	accounts := []Account{}
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.Id, &account.Name); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
