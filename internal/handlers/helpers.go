package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/handlers/pricelist"
	"github.com/macadamiaboy/SigmaPay/internal/postgres"
)

// to fix these structs to make possible to work with this method everywhere
/*
type RequestBody struct {
	EventType pricelist.EventType `json:"event_type"`
}

type Response struct {
	Message    string
	Pricelists []pricelist.EventType
}
*/

// again wrong to use structs out of pricelist package
func HandlerHelper(fn func(*pricelist.RequestBody, *sql.DB) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody pricelist.RequestBody

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db, err := postgres.PrepareDB()
		if err != nil {
			log.Fatalf("failed to prepare the db: %v", err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("Error closing database: %v", err)
			}
		}()

		response, err := fn(&requestBody, db.Connection)
		if err != nil {
			log.Printf("Error during execution: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
