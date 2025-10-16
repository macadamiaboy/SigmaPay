package main

import (
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
)

func main() {
	http.HandleFunc("/pricelists/", handlers.MainHandler(pricelist.GetRequestBody))
	http.HandleFunc("/payments/", handlers.MainHandler(payments.GetRequestBody))
	http.HandleFunc("/addresses/", handlers.MainHandler(addresses.GetRequestBody))
	http.HandleFunc("/events/", handlers.MainHandler(events.GetRequestBody))
	http.HandleFunc("/players/", handlers.MainHandler(players.GetRequestBody))
	http.HandleFunc("/positions/", handlers.MainHandler(positions.GetRequestBody))
	http.HandleFunc("/presence/", handlers.MainHandler(presence.GetRequestBody))
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

	/*
		newEvent := tablesmethods.EventType{Type: "Тренировка", Price: 4500}
		err = newEvent.Save(db.Connection)
		if err != nil {
			log.Println("failed to save")
		} else {
			log.Println("successfully saved")
		}
	*/

	/*
		newEvent, err := pricelist.GetByID(db.Connection, 2)
		if err != nil {
			log.Println("failed to get")
			fmt.Printf("failed to get, err: %v", err)
		} else {
			log.Println("successfully got")
		}
		fmt.Println(newEvent)
	*/

	/*
		newEvent.Delete(db.Connection)
		if err != nil {
			log.Println("failed to delete")
			fmt.Printf("failed to delete, err: %v", err)
		} else {
			log.Println("successfully deleted")
		}
	*/
}
