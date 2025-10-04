package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/postgres"
	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables-methods"
)

func HWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sigma Managing App"))
}

func main() {
	//http.HandleFunc("/", HWorld)
	//log.Fatal(http.ListenAndServe(":8094", nil))

	//init db
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

	newEvent := tablesmethods.EventType{Type: "Тренировка", Price: 4500}
	err = newEvent.Save(db.Connection)
	if err != nil {
		log.Printf("failed to save")
	} else {
		log.Printf("successfully saved")
	}

}
