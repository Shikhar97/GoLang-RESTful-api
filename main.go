package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cozy-software/interview-test/backend/api"
	"github.com/cozy-software/interview-test/backend/internal/database"
)

func init() {
	database.DB = database.New()
}

func main() {
	if len(os.Args[1:]) > 0 && os.Args[1] == "seed" {
		database.Seed(database.DB)
	} else {
		router := api.Mount()
		log.Println("Starting server on port 3000")
		http.ListenAndServe(":3000", router)
	}
}
