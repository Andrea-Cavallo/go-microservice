// main.go
package main

import (
	"myapp/internal/config"
	"myapp/internal/middleware"
	"myapp/internal/router"
	"myapp/internal/utils"
	"net/http"
	"os"
)

func main() {

	log := utils.WithContext().WithField("package", "main")
	log.Infof("Loading mongoClient..")
	// Carica la configurazione e inizializza la connessione a MongoDB
	config.GetMongoClient()
	log.Infof("Configuring zipkin tracer..")
	// Configura il tracer di Zipkin
	tracer := middleware.SetupZipkinTracer()
	log.Infof("Configuring routes..")
	// Configura e avvia il router
	r := router.SetupRouter(tracer)
	log.Infof("Starting server on %s", os.Getenv("SERVICE_IP"))
	log.Fatal(http.ListenAndServe(":8080", r))
}
