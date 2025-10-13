package pricelist

import (
	"encoding/json"
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

func SaveEvent(w http.ResponseWriter, r *http.Request) {
	env := "handlers.payments.SaveEvent"

	var requestBody RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := postgres.PrepareDB()
	if err != nil {
		log.Fatalf("%s: failed to prepare the db: %v", env, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err = requestBody.EventType.Save(db.Connection); err != nil {
		log.Fatalf("%s: failed to save the pricelist: %v", env, err)
	}

	response := Response{
		Message: "New EventType saved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetEventTypeByID(w http.ResponseWriter, r *http.Request) {
	env := "handlers.payments.GetEventTypeByID"
	var eventPrice *pricelist.EventType

	var requestBody RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := postgres.PrepareDB()
	if err != nil {
		log.Fatalf("%s: failed to prepare the db: %v", env, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if eventPrice, err = pricelist.GetByID(db.Connection, int64(requestBody.EventType.Id)); err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	response := Response{
		Message:    "Got pricelist by ID successfully",
		Pricelists: []pricelist.EventType{*eventPrice},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllEventTypes(w http.ResponseWriter, r *http.Request) {
	env := "handlers.payments.GetAllEventTypes"
	var eventPrices *[]pricelist.EventType

	db, err := postgres.PrepareDB()
	if err != nil {
		log.Fatalf("%s: failed to prepare the db: %v", env, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if eventPrices, err = pricelist.GetAll(db.Connection); err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	response := Response{
		Message:    "Got pricelist by ID successfully",
		Pricelists: *eventPrices,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PatchEventType(w http.ResponseWriter, r *http.Request) {
	env := "handlers.payments.PostEventType"

	var requestBody RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := postgres.PrepareDB()
	if err != nil {
		log.Fatalf("%s: failed to prepare the db: %v", env, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	updatedEvent := requestBody.EventType

	if event, err := pricelist.GetByID(db.Connection, int64(updatedEvent.Id)); err != nil {
		log.Println(err)
		http.Error(w, "Event not found", http.StatusNotFound)
	} else {
		if updatedEvent.Type == "" {
			updatedEvent.Type = event.Type
		}
		if updatedEvent.Price == 0 {
			updatedEvent.Price = event.Price
		}
	}

	if err = updatedEvent.Update(db.Connection); err != nil {
		log.Fatalf("%s: failed to update the pricelist: %v", env, err)
	}

	response := Response{
		Message: "Updated pricelist successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	env := "handlers.payments.DeleteEvent"

	var requestBody RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := postgres.PrepareDB()
	if err != nil {
		log.Fatalf("%s: failed to prepare the db: %v", env, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	event, err := pricelist.GetByID(db.Connection, int64(requestBody.EventType.Id))
	if err != nil {
		log.Println(err)
		http.Error(w, "Event not found", http.StatusNotFound)
	}

	if err = event.Delete(db.Connection); err != nil {
		log.Fatalf("%s: failed to save the pricelist: %v", env, err)
	}

	response := Response{
		Message: "EventType record deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
