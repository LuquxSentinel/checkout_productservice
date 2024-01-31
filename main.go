package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	uri := os.Getenv("MONGO_CONN_STR")
	if uri == "" {
		log.Fatal("MONGO_CONN_STR not found in environments")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("MONGO_CONN_STR not found in environments")
	}

	client, err := InitDB(uri)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(dbName).Collection("products")

	storage := NewMongoStorage(collection)
	service := NewService(storage)
	server := NewAPIServer(":5000", service)

	log.Printf("Listening on port : %d", 5000)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
