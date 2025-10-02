package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/postgres"
)

func HWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sigma Managing App"))
}

func main() {
	//http.HandleFunc("/", HWorld)
	//log.Fatal(http.ListenAndServe(":8094", nil))

	if err := postgres.New(); err != nil {
		fmt.Println(err)
		log.Fatal("failed to create")
	}
}
