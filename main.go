package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanaphat2005/simple-bank/api"
	db "github.com/thanaphat2005/simple-bank/db/sqlc"
	"github.com/thanaphat2005/simple-bank/token"
	"github.com/thanaphat2005/simple-bank/utils"
)

const (
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = ":8080"
)

func main() {
	var err error

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)

	}

	ctx := context.Background()
	DB, err := pgxpool.New(ctx, config.DbSource)

	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	defer DB.Close()

	tokenMaker, err := token.NewPasetoMaker(config.SecretKey)
	if err != nil {
		log.Fatal("Cannot create token maker:", err)
	}

	store := db.NewStore(DB)
	server, err := api.NewServer(store, config, tokenMaker)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
