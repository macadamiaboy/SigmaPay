package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/postgres"
)

// to fix these structs to make possible to work with this method everywhere
/*
type RequestBody struct {
	EventType pricelist.EventType `json:"event_type"`
}
*/
type CRUD interface {
	Save(db *sql.DB) error
	Update(db *sql.DB) error
	Delete(db *sql.DB) error
	Get(db *sql.DB) (any, error)
	GetAll(db *sql.DB) (*[]any, error)
}

type Response struct {
	Message string
	Data    *[]any
}

// again wrong to use structs out of pricelist package
func HandlerHelper(bodyGetter func(*http.Request) (CRUD, error), fn func(CRUD, *sql.DB) (*Response, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody CRUD

		requestBody, err := bodyGetter(r)
		if err != nil {
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

		response, err := fn(requestBody, db.Connection)
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
