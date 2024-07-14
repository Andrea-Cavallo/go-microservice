package middleware

import (
	"myapp/internal/utils"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimiterMiddleware limita il numero di richieste che un client può fare in un determinato periodo di tempo
// Burst di 3: Significa che, oltre alle 5 richieste per secondo, il client può fare fino a 3 richieste in più in un colpo solo. Questo buffer di 3 richieste viene utilizzato per gestire picchi momentanei di traffico.
// Quindi, in pratica, un client può fare 8 richieste in un secondo (5 + 3 burst),
// ma se tutte e 8 le richieste vengono fatte immediatamente, il client dovrà attendere per fare ulteriori richieste fino al secondo successivo.
// .
func RateLimiterMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(5, 3) // 5 richiesta al secondo con un burst di 3
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			utils.RespondWithError(w, http.StatusTooManyRequests, "Too Many Requests - slow down my friend")
			return
		}
		next.ServeHTTP(w, r)
	})
}
