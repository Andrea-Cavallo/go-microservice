// middleware/zipkin.go
package middleware

import (
	"myapp/internal/utils"
	"os"

	"github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

// SetupZipkinTracer configura e ritorna un tracer Zipkin
func SetupZipkinTracer() *zipkin.Tracer {
	log := utils.WithContext().WithField("package", "middleware")

	// Ottieni il nome del servizio e l'URL di Zipkin dall'ambiente
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		log.Fatal("SERVICE_NAME environment variable is required")
	}

	zipkinURL := os.Getenv("ZIPKIN_URL")
	if zipkinURL == "" {
		log.Fatal("ZIPKIN_URL environment variable is required")
	}

	// Ottieni l'IP del servizio
	serviceIP := os.Getenv("SERVICE_IP")
	if serviceIP == "" {
		serviceIP = "localhost:8080"
	}

	// Crea un reporter da usare con il tracer.
	reporter := httpreporter.NewReporter(zipkinURL)
	log.WithField("reporter", reporter).WithField("service", serviceName).Info("zipkin")

	// Configura l'endpoint locale per il servizio.
	endpoint, err := zipkin.NewEndpoint(serviceName, serviceIP)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// Configura la strategia di campionamento.
	sampler := zipkin.NewModuloSampler(1)

	// Inizializza il tracer.
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	return tracer
}
