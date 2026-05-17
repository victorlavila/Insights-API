package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"api_insights/internal/controller"
	"api_insights/internal/model"
)

func main() {
	store := model.NewMockStore()
	insightModel := model.NewInsightModel(store)
	insightController := controller.NewInsightController(insightModel)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           insightController.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("api de insights ouvindo em http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
