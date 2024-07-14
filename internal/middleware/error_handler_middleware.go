package middleware

import (
	"myapp/internal/utils"
	"net/http"
)

// ErrorHandlerMiddleware gestisce gli errori in modo centralizzato.
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log := utils.WithContext()
				log.Errorf("Recovered from panic: %v", err)
				utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
