package middleware

import (
	"context"
	"myapp/internal/utils"
	"net/http"
)

type contextKey string

const CorrelationIDKey = contextKey("correlationID")

// CorrelationIDMiddleware genera un ID di correlazione utilizzando la funzione GenerateUUID dal pacchetto utils.
// Viene poi aggiunto agli header della richiesta come "X-Correlation-ID" e agli header della risposta.
// Il middleware poi chiama il prossimo handler nella catena.
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID, err := utils.GenerateUUID()
		if err != nil {
			http.Error(w, "Unable to generate correlation ID", http.StatusInternalServerError)
			return
		}

		// Aggiungi il correlation ID agli header della richiesta e della risposta
		r.Header.Set("X-Correlation-ID", correlationID)
		w.Header().Set("X-Correlation-ID", correlationID)

		// Aggiungi il correlation ID al contesto della richiesta
		ctx := context.WithValue(r.Context(), CorrelationIDKey, correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCorrelationID recupera il correlation ID dal contesto
func GetCorrelationID(ctx context.Context) string {
	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return correlationID
	}
	return ""
}
