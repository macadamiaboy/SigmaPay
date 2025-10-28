package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	"github.com/macadamiaboy/SigmaPay/internal/handlers/addresses"
	"github.com/macadamiaboy/SigmaPay/internal/handlers/events"
	"github.com/macadamiaboy/SigmaPay/internal/handlers/payments"
	"github.com/macadamiaboy/SigmaPay/internal/handlers/players"
	"github.com/macadamiaboy/SigmaPay/internal/handlers/positions"
	"github.com/macadamiaboy/SigmaPay/internal/handlers/presence"
	"github.com/macadamiaboy/SigmaPay/internal/handlers/pricelist"
	"github.com/macadamiaboy/SigmaPay/internal/postgres"
)

func main() {
	if err := postgres.New(); err != nil {
		fmt.Println(err)
		log.Fatal("failed to create")
	}

	http.HandleFunc("/pricelists/", handlers.CRUDHandler(pricelist.GetRequestBody))

	http.HandleFunc("/payments/", handlers.CRUDHandler(payments.GetRequestBody))
	http.HandleFunc("/payments/debts", payments.DebtHandler)

	http.HandleFunc("/addresses/", handlers.CRUDHandler(addresses.GetRequestBody))

	http.HandleFunc("/events/", handlers.CRUDHandler(events.GetRequestBody))
	http.HandleFunc("/events/type/", events.ByTypeHandler)

	http.HandleFunc("/players/", handlers.CRUDHandler(players.GetRequestBody))

	http.HandleFunc("/positions/", handlers.CRUDHandler(positions.GetRequestBody))

	http.HandleFunc("/presence/", handlers.CRUDHandler(presence.GetRequestBody))

	log.Fatal(http.ListenAndServe(":8094", nil))

	//init db
	/*
		if err := postgres.New(); err != nil {
			fmt.Println(err)
			log.Fatal("failed to create")
		}

		db, err := postgres.PrepareDB()
		if err != nil {
			log.Fatal("failed to init")
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("Error closing database: %v", err)
			}
		}()
	*/

}
