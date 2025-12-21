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

	cfg := config.LoadDBConfigData()
	addr := fmt.Sprintf("%s:%v", cfg.Server.Host, cfg.Server.Port)

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	// routes for positions, not so necessary
	// needed just to sort players by positions
	// think about creating the list with the initiation of the db
	router.Route("/positions", func(r chi.Router) {
		requestBody := positions.GetRequestBody

		r.Get("/", handlers.CRUDHandler(requestBody, handlers.GetAllHelper))
		r.Post("/", handlers.CRUDHandler(requestBody, handlers.SaveHelper))
		r.Delete("/", handlers.CRUDHandler(requestBody, handlers.DeleteHelper))
	})

	// the same with the pricelist: no need to post anything, just to store the prices
	// maybe create it with the init also and just make possible to update
	router.Route("/pricelist", func(r chi.Router) {
		requestBody := pricelist.GetRequestBody

		r.Get("/", handlers.CRUDHandler(requestBody, handlers.GetAllHelper))
		r.Post("/", handlers.CRUDHandler(requestBody, handlers.GetHelper))
		r.Delete("/", handlers.CRUDHandler(requestBody, handlers.GetHelper))
		r.Patch("/", handlers.CRUDHandler(requestBody, handlers.PatchHelper))
	})

	// the same again: don't need much actions, have a lot connected with just two
	// if needed, maybe add addresses for away games, but two basic create with the init
	router.Route("/addresses", func(r chi.Router) {
		requestBody := addresses.GetRequestBody

		r.Get("/", handlers.CRUDHandler(requestBody, handlers.GetHelper))
		r.Post("/", handlers.CRUDHandler(requestBody, handlers.SaveHelper))
		r.Delete("/", handlers.CRUDHandler(requestBody, handlers.DeleteHelper))
	})

	// there's ByTypeHandler. Check if it's necessary and add a route
	router.Route("/events", func(r chi.Router) {
		requestBody := events.GetRequestBody

		r.Get("/", handlers.CRUDHandler(requestBody, handlers.GetHelper))
		r.Post("/", handlers.CRUDHandler(requestBody, handlers.SaveHelper))
		r.Delete("/", handlers.CRUDHandler(requestBody, handlers.DeleteHelper))
		r.Patch("/", handlers.CRUDHandler(requestBody, handlers.PatchHelper))

		r.Get("/all", handlers.CRUDHandler(requestBody, handlers.GetAllHelper))
		r.Get("/payments", events.PaymentsHandler)
	})

	router.Route("/players", func(r chi.Router) {
		requestBody := players.GetRequestBody

		r.Get("/", handlers.CRUDHandler(requestBody, handlers.GetHelper))
		r.Post("/", handlers.CRUDHandler(requestBody, handlers.SaveHelper))
		r.Delete("/", handlers.CRUDHandler(requestBody, handlers.DeleteHelper))
		r.Patch("/", handlers.CRUDHandler(requestBody, handlers.PatchHelper))
	})

	router.Route("/presence", func(r chi.Router) {
		requestBody := presence.GetRequestBody

		r.Post("/", handlers.CRUDHandler(requestBody, handlers.SaveHelper))
		r.Delete("/", handlers.CRUDHandler(requestBody, handlers.DeleteHelper))
	})

	router.Route("/payments", func(r chi.Router) {
		requestBody := payments.GetRequestBody

		r.Get("/", handlers.CRUDHandler(requestBody, handlers.GetHelper))
		r.Post("/", handlers.CRUDHandler(requestBody, handlers.SaveHelper))
		r.Delete("/", handlers.CRUDHandler(requestBody, handlers.DeleteHelper))
		r.Patch("/", handlers.CRUDHandler(requestBody, handlers.PatchHelper))

		r.Get("/all", handlers.CRUDHandler(requestBody, handlers.GetAllHelper))
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
