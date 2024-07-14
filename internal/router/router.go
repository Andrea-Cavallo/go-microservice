package router

import (
	"myapp/internal/handlers"
	"myapp/internal/middleware"
	"myapp/internal/utils"
	"myapp/internal/utils/constants"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	zipkinMiddleware "github.com/openzipkin/zipkin-go/middleware/http"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRouter configura le rotte HTTP per l'applicazione.
func SetupRouter(tracer *zipkin.Tracer) *mux.Router {
	// Crea un nuovo router
	r := mux.NewRouter()

	// Crea il middleware di Zipkin per il server
	zipkinMw := zipkinMiddleware.NewServerMiddleware(tracer, zipkinMiddleware.TagResponseSize(true))
	r.Use(zipkinMw)

	// Applica middleware globali

	r.Use(middleware.RateLimiterMiddleware)
	r.Use(middleware.CorrelationIDMiddleware)
	r.Use(middleware.ErrorHandlerMiddleware)

	// Definizione rotta per gli utenti
	userRoutes := r.PathPrefix(constants.USERS).Subrouter()
	userRoutes.HandleFunc(constants.BLANK, handlers.GetUsers(tracer)).Methods(constants.HTTPGet)
	userRoutes.HandleFunc(constants.BLANK, handlers.CreateUser(tracer)).Methods(constants.HTTPPost)
	userRoutes.HandleFunc(constants.ID, handlers.GetUserByID(tracer)).Methods(constants.HTTPGet)
	userRoutes.HandleFunc(constants.ID, handlers.DeleteUserByID(tracer)).Methods(constants.HTTPDelete)
	userRoutes.HandleFunc(constants.ID, handlers.UpdateUser(tracer)).Methods(constants.HTTPPut)

	// Aggiunge una rotta per le metriche di Prometheus
	r.Handle("/metrics", handlers.MetricsHandler())

	// Aggiunge una rotta per la documentazione Swagger
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Aggiunge un handler per gestire le rotte non trovate (404)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return r
}

// notFoundHandler gestisce gli errori 404 per le rotte non definite.
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusNotFound, "Rotta non mappata, sconosciuta")
}
