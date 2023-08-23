package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"stg-go-websocket-server/api"
	"stg-go-websocket-server/util"
	"stg-go-websocket-server/ws"
)

func main() {

	config, err := util.LoadConfig(".")
	log.Printf("Config: %v", config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://jleva:gmjnhoHk2vofz2S8@cabbage-order.ilyal.mongodb.net/slide-show?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("stg-chat").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}

	s, err := api.NewServer(config, ws.NewManager(client))
	if err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("failed to create server %v", err))
	}

	err = s.Start()
	if err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("failed to start server %v", err))
	}
}
