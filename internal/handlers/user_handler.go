package handlers

import (
	"encoding/json"
	"myapp/internal/middleware"
	"myapp/internal/models"
	"myapp/internal/services"
	"myapp/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
)

//NOTA, l'uso di una funzione che ritorna un'altra funzione è un pattern comune - Closure ( spesso usato negli HTTP handlers )
//CLOSURE PATTERN Le closure permettono di catturare e mantenere il contesto delle variabili presenti al momento della loro definizione.
//In questo caso specifico, la closure cattura il tracer, che viene poi utilizzato nell'handler HTTP per tracciare le richieste.
//Puoi facilmente cambiare il comportamento dell'handler passando un diverso tracer.
//La closure mantiene il contesto del tracer al momento della definizione, garantendo che lo stesso tracer venga utilizzato per ogni richiesta gestita dall'handler.

// GetUsers recupera tutti gli utenti e li restituisce come risposta JSON.
// @Summary Get all users
// @Description Recupera tutti gli utenti
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(tracer *zipkin.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := utils.WithContext()

		correlationID := middleware.GetCorrelationID(r.Context())
		log.Infof("GetAllUsers Handler with - correlationID: %s", correlationID)

		// Crea uno span per tracciare l'operazione GetUsers
		span := tracer.StartSpan("GetUsers")
		defer span.Finish()

		// Crea un nuovo contesto con lo span corrente
		ctx := zipkin.NewContext(r.Context(), span)

		// Passa lo span e il contesto al servizio
		users, err := services.GetAllUsers(ctx, tracer)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving users")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, users)
	}
}

// CreateUser decodifica il JSON in ingresso dalla richiesta e crea un nuovo utente.
// @Summary Create a new user
// @Description Crea un nuovo utente
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body  models.User  true  "User object"
// @Success 201 {object} models.User
// @Router /users [post]
func CreateUser(tracer *zipkin.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := utils.WithContext()

		correlationID := middleware.GetCorrelationID(r.Context())
		log.Info("CreateUsers Handler with - correlationID: %s", correlationID)

		// Crea uno span per tracciare l'operazione CreateUser
		span := tracer.StartSpan("CreateUser")
		defer span.Finish()

		// Assicura che il corpo della richiesta venga chiuso alla fine della funzione
		defer utils.CloseRequestBody(r.Body)

		// Crea una variabile per memorizzare i dati dell'utente decodificati
		var user models.User

		// Tenta di decodificare il JSON nel corpo della richiesta nella variabile user
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Crea un nuovo utente tramite il servizio
		createdUser, err := services.CreateUser(user)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
			return
		}
		utils.RespondWithJSON(w, http.StatusCreated, createdUser)
	}
}

// GetUserByID recupera un utente per ID e risponde con il JSON dell'utente se trovato, altrimenti con un messaggio di errore.
// @Summary Get user by ID
// @Description Recupera un utente per ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} utils.Response
// @Router /users/{id} [get]
func GetUserByID(tracer *zipkin.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := utils.WithContext()

		correlationID := middleware.GetCorrelationID(r.Context())
		log.Info("GetUserByID Handler with - correlationID: %s", correlationID)

		// Crea uno span per tracciare l'operazione GetUserByID
		span := tracer.StartSpan("GetUserByID")
		defer span.Finish()

		// Ottiene i parametri della route dalla richiesta
		params := mux.Vars(r)

		// Recupera il valore del parametro "id" dalla mappa
		id := params["id"]

		// Utilizza l'ID per recuperare l'utente corrispondente
		user, err := services.GetUserByID(id)
		if err != nil {
			// Se l'utente non viene trovato, risponde con un errore 404 (Not Found)
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// Se l'utente viene trovato, risponde con i dati dell'utente
		utils.RespondWithJSON(w, http.StatusOK, user)
	}
}

// DeleteUserByID elimina un utente per ID.
// @Summary Delete a user by ID
// @Description Elimina un utente per ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} utils.Response
// @Router /users/{id} [delete]
func DeleteUserByID(tracer *zipkin.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := utils.WithContext()

		correlationID := middleware.GetCorrelationID(r.Context())
		log.Info("DeleteUserById Handler with - correlationID: %s", correlationID)
		// Crea uno span per tracciare l'operazione DeleteUserByID
		span := tracer.StartSpan("DeleteUserByID")
		defer span.Finish()

		// Ottiene i parametri della route dalla richiesta
		params := mux.Vars(r)

		// Elimina l'utente tramite il servizio
		err := services.DeleteUserByID(params["id"])
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Error deleting user")
			return
		}
		utils.RespondWithJSON(w, http.StatusNoContent, nil)
	}
}

// UpdateUser aggiorna un utente in base al payload della richiesta.
// @Summary Update a user by ID
// @Description Aggiorna un utente per ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "User ID"
// @Param   user  body  models.User  true  "User object"
// @Success 200 {object} models.User
// @Failure 404 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /users/{id} [put]
func UpdateUser(tracer *zipkin.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := utils.WithContext()

		correlationID := middleware.GetCorrelationID(r.Context())
		log.Info("UpdateUser Handler with - correlationID: %s", correlationID)
		// Crea uno span per tracciare l'operazione UpdateUser
		span := tracer.StartSpan("UpdateUser")
		defer span.Finish()

		// Assicura che il corpo della richiesta venga chiuso alla fine della funzione
		defer utils.CloseRequestBody(r.Body)

		// Crea una variabile per memorizzare i dati dell'utente decodificati
		var user models.User

		// Tenta di decodificare il JSON nel corpo della richiesta nella variabile user
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Ottiene i parametri della route dalla richiesta
		params := mux.Vars(r)

		// Aggiorna l'utente tramite il servizio
		updatedUser, err := services.UpdateUser(params["id"], user)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Error updating user")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, updatedUser)
	}
}
