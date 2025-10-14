package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func SaveHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.SaveHelper"

	if err := requestBody.Save(db); err != nil {
		log.Fatalf("%s: failed to save the pricelist: %v", env, err)
	}

	response := Response{
		Message: "New EventType saved successfully",
	}

	return &response, nil
}

func GetHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.GetHelper"

	eventPrice, err := requestBody.Get(db)
	if err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	response := Response{
		Message: "Got pricelist by ID successfully",
		Data:    &[]any{eventPrice},
	}

	return &response, nil
}

func GetAllHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.GetAllHelper"

	eventPrices, err := requestBody.GetAll(db)
	if err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	response := Response{
		Message: "Got all pricelists successfully",
		Data:    eventPrices,
	}

	return &response, nil
}

func DeleteHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.DeleteHelper"

	event, err := requestBody.Get(db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("not found")
		//http.Error(w, "Event not found", http.StatusNotFound)
	}

	if et, ok := event.(CRUD); ok {
		if err = et.Delete(db); err != nil {
			log.Fatalf("%s: failed to save the pricelist: %v", env, err)
		}
	}

	response := Response{
		Message: "EventType record deleted successfully",
	}

	return &response, nil
}

func PatchHelper(requestBody CRUD, db *sql.DB) (*Response, error) {
	env := "handlers.helpers.PatchHelper"

	if err := requestBody.Update(db); err != nil {
		log.Fatalf("%s: failed to update the pricelist: %v", env, err)
	}

	response := Response{
		Message: "Updated pricelist successfully",
	}

	return &response, nil
}
