package pricelist

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/postgres"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/pricelist"
)

type RequestBody struct {
	Id int64 `json:"id"`
}

type Response struct {
	Message    string
	Pricelists []pricelist.EventType
}

func GetPricelist(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			GetAllPricelists(w, r)
			return
		}
	*/

	/*
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
	*/

	//GetPricelistByID(w, r, requestBody.Id)
}

func GetPricelistByID(w http.ResponseWriter, r *http.Request) {
	env := "handlers.payments.GetPricelistByID"
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

	if eventPrice, err = pricelist.GetByID(db.Connection, int64(requestBody.Id)); err != nil {
		log.Fatalf("%s: failed to get the pricelist: %v", env, err)
	}

	/*
		jsonData, err := json.Marshal(eventPrice)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	*/

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

func GetAllPricelists(w http.ResponseWriter, r *http.Request) {
	env := "handlers.payments.GetPricelist"
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
