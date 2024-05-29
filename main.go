package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kunal768/bitespeed/client"
	"github.com/kunal768/bitespeed/parser"
)

func initDbClient() *pgxpool.Pool {
	url := os.Getenv("DATABASE_URL")
	client, err := client.ConnectDB(url)
	if err != nil {
		panic(err)
	}
	return client
}

func initParserService(client *pgxpool.Pool) parser.Service {
	repo := parser.NewRepository(client)
	return parser.NewParserService(repo)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	dbClient := initDbClient()
	parserService := initParserService(dbClient)
	router := gin.Default()

	router.POST("/identify", parser.HandleContactRequest(parserService))

	router.Run(":4000")
}
