package events

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	"github.com/macadamiaboy/SigmaPay/internal/postgres"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/events"
)

type RequestBody struct {
	Event events.Event `json:"event"`
}

func GetRequestBody(r *http.Request) (handlers.CRUD, error) {
	var requestBody *RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("failed to get the request body, err: %w", err)
	}

	return &requestBody.Event, nil
}

func ByMonthHandler(db *postgres.DataBase /*, reqMonth int, reqYear int*/) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "There's no such method", http.StatusMethodNotAllowed)
			return
		}

		requestBody, err := GetRequestBody(r)
		if err != nil {
			log.Fatalf("failed to get the request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event, ok := requestBody.(*events.Event)
		if !ok {
			msg := "provided request body is not of the expected type"
			log.Fatal(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		reqYear := chi.URLParam(r, "year")
		reqMonth := chi.URLParam(r, "month")

		data, err := event.GetAllByMonth(db.Connection, reqMonth, reqYear)
		if err != nil {
			log.Fatalf("failed to get the data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := handlers.Response{
			Message: "Got all records successfully",
			Data:    data,
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ByTypeHandler(db *postgres.DataBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fn func(*sql.DB) (*[]any, error)

		if r.Method != http.MethodGet {
			http.Error(w, "There's no such method", http.StatusMethodNotAllowed)
			return
		}

		requestBody, err := GetRequestBody(r)
		if err != nil {
			log.Fatalf("failed to get the request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event, ok := requestBody.(*events.Event)
		if !ok {
			msg := "provided request body is not of the expected type"
			log.Fatal(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		fn = event.GetAllGames

		data, err := fn(db.Connection)
		if err != nil {
			log.Fatalf("failed to get the data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := handlers.Response{
			Message: "Got all records successfully",
			Data:    data,
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PaymentsHandler(db *postgres.DataBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := GetRequestBody(r)
		if err != nil {
			log.Fatalf("failed to get the request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event, ok := requestBody.(*events.Event)
		if !ok {
			msg := "provided request body is not of the expected type"
			log.Fatal(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		data, err := event.GetAllEventPayments(db.Connection)
		if err != nil {
			log.Fatalf("failed to get the data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := handlers.Response{
			Message: "Got all records successfully",
			Data:    data,
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
