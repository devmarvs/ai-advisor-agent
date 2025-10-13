package main

import (
	"log"
	"os"
)

func main() {
	r := SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	log.Printf("API listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
