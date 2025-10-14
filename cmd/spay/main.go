package main

import (
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	pricelist "github.com/macadamiaboy/SigmaPay/internal/handlers/pricelist"
)

func main() {
	http.HandleFunc("/pricelists/all", handlers.HandlerHelper(pricelist.GetRequestBody, handlers.GetAllHelper))
	http.HandleFunc("/pricelists/", handlers.HandlerHelper(pricelist.GetRequestBody, handlers.GetHelper))
	http.HandleFunc("/pricelists/save", handlers.HandlerHelper(pricelist.GetRequestBody, handlers.SaveHelper))
	http.HandleFunc("/pricelists/delete", handlers.HandlerHelper(pricelist.GetRequestBody, handlers.DeleteHelper))
	http.HandleFunc("/pricelists/update", handlers.HandlerHelper(pricelist.GetRequestBody, handlers.PatchHelper))
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
