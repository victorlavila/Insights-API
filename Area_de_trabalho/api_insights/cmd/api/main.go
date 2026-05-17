package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"api_insights/internal/repository"
	"api_insights/internal/service"
	httptransport "api_insights/internal/transport/http"
)

func main() {
	store := repository.NewMockStore()
	insightService := service.NewInsightService(store)
	handler := httptransport.NewHandler(insightService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("api de insights ouvindo em http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
