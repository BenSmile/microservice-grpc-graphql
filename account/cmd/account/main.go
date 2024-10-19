package main

import (
	"log"
	"time"

	"github.com/bensmile/microservice-grpc-graphql/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
}

func main() {

	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		r, err = account.NewPostresRepository(cfg.DatabaseUrl)
		log.Println(err)
		if err != nil {
			return err
		}
		return nil
	})

	defer r.Close()

	log.Println("Listening on port 8080 ...")

	s := account.NewService(r)

	if err := account.ListenGrpc(s, 8080); err != nil {
		log.Fatal(err)
	}

}
