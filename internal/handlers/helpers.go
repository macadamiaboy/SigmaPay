package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/macadamiaboy/SigmaPay/internal/postgres"
)

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

func MainHandler(bodyGetter func(*http.Request) (CRUD, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fn func(CRUD, *sql.DB) (*Response, error)

		switch r.Method {
		case http.MethodGet:
			if strings.Contains(r.URL.Path, "/all") {
				fn = GetAllHelper
			} else {
				fn = GetHelper
			}
		case http.MethodPost:
			fn = SaveHelper
		case http.MethodPatch:
			fn = PatchHelper
		case http.MethodDelete:
			fn = DeleteHelper
		default:
			http.Error(w, "There's no such method", http.StatusMethodNotAllowed)
			return
		}

		requestBody, err := bodyGetter(r)
		if err != nil {
			log.Fatalf("failed to get the request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db, err := postgres.PrepareDB()
		if err != nil {
			log.Fatalf("failed to prepare the db: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("Error closing database: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		response, err := fn(requestBody, db.Connection)
		if err != nil {
			log.Printf("Error during execution: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func SaveHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.SaveHelper"

	if err := requestBody.Save(db); err != nil {
		log.Fatalf("%s: failed to save the record: %v", env, err)
		return nil, fmt.Errorf("failed to save the record: %w", err)
	}

	response := Response{
		Message: "New record saved successfully",
	}

	return &response, nil
}

func GetHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.GetHelper"

	eventPrice, err := requestBody.Get(db)
	if err != nil {
		log.Fatalf("%s: failed to get the record: %v", env, err)
		return nil, fmt.Errorf("failed to get the record: %w", err)
	}

	response := Response{
		Message: "Got record by ID successfully",
		Data:    &[]any{eventPrice},
	}

	return &response, nil
}

func GetAllHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.GetAllHelper"

	eventPrices, err := requestBody.GetAll(db)
	if err != nil {
		log.Fatalf("%s: failed to get records: %v", env, err)
		return nil, fmt.Errorf("failed to get records: %w", err)
	}

	response := Response{
		Message: "Got all records successfully",
		Data:    eventPrices,
	}

	return &response, nil
}

func DeleteHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.DeleteHelper"

	event, err := requestBody.Get(db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("the record is not found: %w", err)
	}

	if et, ok := event.(CRUD); ok {
		if err = et.Delete(db); err != nil {
			log.Fatalf("%s: failed to delete the record: %v", env, err)
			return nil, fmt.Errorf("failed to delete the record: %w", err)
		}
	}

	response := Response{
		Message: "Record deleted successfully",
	}

	return &response, nil
}

func PatchHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.PatchHelper"

	if err := requestBody.Update(db); err != nil {
		log.Fatalf("%s: failed to update the record: %v", env, err)
		return nil, fmt.Errorf("failed to update the record: %w", err)
	}

	response := Response{
		Message: "Updated record successfully",
	}

	return &response, nil
}
