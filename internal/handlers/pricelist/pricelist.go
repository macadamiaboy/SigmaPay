package pricelist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/postgres"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/pricelist"
)

type RequestBody struct {
	EventType pricelist.EventType `json:"event_type"`
}

type Response struct {
	Message    string
	Pricelists []pricelist.EventType
}

func HandlerHelper(fn func(*RequestBody, *sql.DB) (*Response, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody RequestBody

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

func SaveEvent(requestBody *RequestBody, db *sql.DB) (*Response, error) {
	env := "handlers.payments.SaveEvent"

	if err := requestBody.EventType.Save(db); err != nil {
		log.Fatalf("%s: failed to save the pricelist: %v", env, err)
	}

	response := Response{
		Message: "New EventType saved successfully",
	}

	return &response, nil
}

func GetEventTypeByID(requestBody *RequestBody, db *sql.DB) (*Response, error) {
	env := "handlers.payments.GetEventTypeByID"

	var eventPrice *pricelist.EventType

	eventPrice, err := pricelist.GetByID(db, int64(requestBody.EventType.Id))
	if err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	response := Response{
		Message:    "Got pricelist by ID successfully",
		Pricelists: []pricelist.EventType{*eventPrice},
	}

	return &response, nil
}

func GetAllEventTypes(requestBody *RequestBody, db *sql.DB) (*Response, error) {
	env := "handlers.payments.GetAllEventTypes"

	var eventPrices *[]pricelist.EventType

	eventPrices, err := pricelist.GetAll(db)
	if err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	response := Response{
		Message:    "Got all pricelists successfully",
		Pricelists: *eventPrices,
	}

	return &response, nil
}

func PatchEventType(requestBody *RequestBody, db *sql.DB) (*Response, error) {
	env := "handlers.payments.PostEventType"

	updatedEvent := requestBody.EventType

	if event, err := pricelist.GetByID(db, int64(updatedEvent.Id)); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("not found")
		//http.Error(w, "Event not found", http.StatusNotFound)
	} else {
		if updatedEvent.Type == "" {
			updatedEvent.Type = event.Type
		}
		if updatedEvent.Price == 0 {
			updatedEvent.Price = event.Price
		}
	}

	if err := updatedEvent.Update(db); err != nil {
		log.Fatalf("%s: failed to update the pricelist: %v", env, err)
	}

	response := Response{
		Message: "Updated pricelist successfully",
	}

	return &response, nil
}

func DeleteEvent(requestBody *RequestBody, db *sql.DB) (*Response, error) {
	env := "handlers.payments.DeleteEvent"

	event, err := pricelist.GetByID(db, int64(requestBody.EventType.Id))
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("not found")
		//http.Error(w, "Event not found", http.StatusNotFound)
	}

	if err = event.Delete(db); err != nil {
		log.Fatalf("%s: failed to save the pricelist: %v", env, err)
	}

	response := Response{
		Message: "EventType record deleted successfully",
	}

	return &response, nil
}
