package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file:", err)
	}

	store, err := NewPostgresStore()
	if err != nil { // issue with creating our postgresstore
		log.Fatal(err)
	}
	defer store.db.Close() // close the db after we exit (from an error or something else)

	if err := store.Setup(); err != nil { // issue w/ setup (i.e. table creation failed)
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	server.Start()
}
