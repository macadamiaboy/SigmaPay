package players

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func DebtHandler(w http.ResponseWriter, r *http.Request) {
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

	event, ok := requestBody.(*players.Player)
	if !ok {
		msg := "provided request body is not of the expected type"
		log.Fatal(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if strings.Contains(r.URL.Path, "/total") {
		fn = event.GetTotalDebt
	} else if strings.Contains(r.URL.Path, "/unpayed") {
		fn = event.GetAllPlayersDebts
	} else if strings.Contains(r.URL.Path, "/all") {
		fn = event.GetAllPlayersPayments
	} else {
		http.Error(w, "There's no such method", http.StatusMethodNotAllowed)
		return
	}

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
