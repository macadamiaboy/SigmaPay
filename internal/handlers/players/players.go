package players

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	"github.com/macadamiaboy/SigmaPay/internal/postgres"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/players"
)

type RequestBody struct {
	Player players.Player `json:"player"`
}

func GetRequestBody(r *http.Request) (handlers.CRUD, error) {
	var requestBody *RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("failed to get the request body, err: %w", err)
	}

	return &requestBody.Player, nil
}

func debtHelper(db *postgres.DataBase, fn func(*sql.DB, *players.Player) (*[]any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := GetRequestBody(r)
		if err != nil {
			log.Fatalf("failed to get the request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		player, ok := requestBody.(*players.Player)
		if !ok {
			msg := "provided request body is not of the expected type"
			log.Fatal(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		data, err := fn(db.Connection, player)
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

func GetAllPlayersPaymentsHandler(db *postgres.DataBase) http.HandlerFunc {
	return debtHelper(db, players.GetAllPlayersPayments)
}

func GetAllPlayersDebtsHandler(db *postgres.DataBase) http.HandlerFunc {
	return debtHelper(db, players.GetAllPlayersDebts)
}

func GetTotalDebtHandler(db *postgres.DataBase) http.HandlerFunc {
	return debtHelper(db, players.GetTotalDebt)
}

func SigmaHandler(db *postgres.DataBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := GetRequestBody(r)
		if err != nil {
			log.Fatalf("failed to get the request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		player, ok := requestBody.(*players.Player)
		if !ok {
			msg := "provided request body is not of the expected type"
			log.Fatal(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		data, err := player.GetAllSigma(db.Connection)
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
