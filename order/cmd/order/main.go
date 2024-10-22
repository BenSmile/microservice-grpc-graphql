package main

import (
	"log"
	"time"

	"github.com/bensmile/microservice-grpc-graphql/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
	AccountUrl  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogUrl  string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		r, err = order.NewPostgresRepository(cfg.DatabaseUrl)
		if err != nil {
			log.Println(err)
			return err
		}
		defer r.Close()
		return nil
	})
	log.Println("Listening on port 8080...")
	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountUrl, cfg.CatalogUrl, 8080))

}
