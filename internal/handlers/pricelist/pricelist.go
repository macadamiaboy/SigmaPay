package pricelist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/pricelist"
)

type RequestBody struct {
	EventType pricelist.EventType `json:"event_type"`
}

func GetRequestBody(r *http.Request) (handlers.CRUD, error) {
	var requestBody *RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("failed to get the request body, err: %w", err)
	}

	return &requestBody.EventType, nil
}

func SaveEventType(requestBody handlers.CRUD, db *sql.DB) (*handlers.Response, error) {
	env := "handlers.payments.SaveEvent"

	if err := requestBody.Save(db); err != nil {
		log.Fatalf("%s: failed to save the pricelist: %v", env, err)
	}

	response := handlers.Response{
		Message: "New EventType saved successfully",
	}

	return &response, nil
}

func GetEventTypeByID(requestBody handlers.CRUD, db *sql.DB) (*handlers.Response, error) {
	env := "handlers.payments.GetEventTypeByID"

	eventPrice, err := requestBody.Get(db)
	if err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	response := handlers.Response{
		Message: "Got pricelist by ID successfully",
		Data:    &[]any{eventPrice},
	}

	return &response, nil
}

func GetAllEventTypes(requestBody handlers.CRUD, db *sql.DB) (*handlers.Response, error) {
	env := "handlers.payments.GetAllEventTypes"

	eventPrices, err := requestBody.GetAll(db)
	if err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	response := handlers.Response{
		Message: "Got all pricelists successfully",
		Data:    eventPrices,
	}

	return &response, nil
}

func PatchEventType(requestBody handlers.CRUD, db *sql.DB) (*handlers.Response, error) {
	env := "handlers.payments.PostEventType"

	updatedEvent := requestBody

	event, err := requestBody.Get(db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("not found")
		//http.Error(w, "Event not found", http.StatusNotFound)
	}

	checkedEvent, evChecked := event.(*pricelist.EventType)
	checkedUpdEvent, updEvChecked := updatedEvent.(*pricelist.EventType)

	if evChecked && updEvChecked {
		if checkedUpdEvent.Type == "" {
			checkedUpdEvent.Type = checkedEvent.Type
		}
		if checkedUpdEvent.Price == 0 {
			checkedUpdEvent.Price = checkedEvent.Price
		}
	}

	if err := updatedEvent.Update(db); err != nil {
		log.Fatalf("%s: failed to update the pricelist: %v", env, err)
	}

	response := handlers.Response{
		Message: "Updated pricelist successfully",
	}

	return &response, nil
}

func DeleteEventType(requestBody handlers.CRUD, db *sql.DB) (*handlers.Response, error) {
	env := "handlers.payments.DeleteEvent"

	event, err := requestBody.Get(db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("not found")
		//http.Error(w, "Event not found", http.StatusNotFound)
	}

	if et, ok := event.(*pricelist.EventType); ok {
		if err = et.Delete(db); err != nil {
			log.Fatalf("%s: failed to save the pricelist: %v", env, err)
		}
	}

	response := handlers.Response{
		Message: "EventType record deleted successfully",
	}

	return &response, nil
}
