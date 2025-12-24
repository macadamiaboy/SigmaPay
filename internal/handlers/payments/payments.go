package payments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	"github.com/macadamiaboy/SigmaPay/internal/postgres"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/payments"
)

type RequestBody struct {
	Payment payments.Payment `json:"payment"`
}

func GetRequestBody(r *http.Request) (handlers.CRUD, error) {
	var requestBody *RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("failed to get the request body, err: %w", err)
	}

	return &requestBody.Payment, nil
}

func DebtHandler(db *postgres.DataBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := GetRequestBody(r)
		if err != nil {
			log.Fatalf("failed to get the request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payment, ok := requestBody.(*payments.Payment)
		if !ok {
			msg := "provided request body is not of the expected type"
			log.Fatal(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		data, err := payment.GetDebts(db.Connection)
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
