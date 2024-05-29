package main

import (
	"log"
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
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.POST("/identify", parser.HandleContactRequest(parserService))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

}
