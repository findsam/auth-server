package main

import (
	"fmt"
	"log"

	"github.com/findsam/food-server/api"
	"github.com/findsam/food-server/config"
	"github.com/findsam/food-server/db"
)

func main() {
	mongoClient, err := db.ConnectToMongo(config.Envs.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), mongoClient)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
