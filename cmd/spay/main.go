package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/macadamiaboy/SigmaPay/internal/config"
	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	"github.com/macadamiaboy/SigmaPay/internal/handlers/players"
	"github.com/macadamiaboy/SigmaPay/internal/postgres"
)

func main() {
	if err := postgres.New(); err != nil {
		fmt.Println(err)
		log.Fatal("failed to create")
	}

	/*
		http.HandleFunc("/pricelists/", handlers.CRUDHandler(pricelist.GetRequestBody))

		http.HandleFunc("/payments/", handlers.CRUDHandler(payments.GetRequestBody))
		http.HandleFunc("/payments/debts", payments.DebtHandler)

		http.HandleFunc("/addresses/", handlers.CRUDHandler(addresses.GetRequestBody))

		http.HandleFunc("/events/", handlers.CRUDHandler(events.GetRequestBody))
		http.HandleFunc("/events/type/", events.ByTypeHandler)

		http.HandleFunc("/players/", handlers.CRUDHandler(players.GetRequestBody))
		http.HandleFunc("/players/debts/", players.DebtHandler)

		http.HandleFunc("/positions/", handlers.CRUDHandler(positions.GetRequestBody))

		http.HandleFunc("/presence/", handlers.CRUDHandler(presence.GetRequestBody))

		log.Fatal(http.ListenAndServe(":8094", nil))
	*/

	cfg := config.LoadDBConfigData()
	addr := fmt.Sprintf("%s:%v", cfg.Server.Host, cfg.Server.Port)

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	router.Get("/players/", handlers.CRUDHandler(players.GetRequestBody, handlers.GetHelper))
	router.Get("/hello/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	router.Route("/pricelists", func(r chi.Router) {

	})

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	log.Printf("starting server. address: %s", addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("failed to start the server")
	}

	log.Fatal("server stoppped")

}
