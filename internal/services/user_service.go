package services

import (
	"context"
	"errors"
	"myapp/internal/models"
	"myapp/internal/repository"
	"myapp/internal/utils"

	"github.com/openzipkin/zipkin-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllUsers retrieves all users from the MongoDB collection
func GetAllUsers(ctx context.Context, tracer *zipkin.Tracer) ([]models.User, error) {
	log := utils.WithContext()
	log.Info("Get all users...")

	// Recupera lo span dal contesto
	parentSpan := zipkin.SpanFromContext(ctx)
	// Crea uno span figlio
	childSpan := tracer.StartSpan("GetAllUsers_Service", zipkin.Parent(parentSpan.Context()))
	defer childSpan.Finish()

	// Crea un nuovo contesto con lo span figlio
	childCtx := zipkin.NewContext(ctx, childSpan)

	users, err := repository.GetUsers(childCtx, tracer)
	if err != nil {
		log.Errorf("Errore durante la getAll: %v", err)
	}
	return users, err
}

// CreateUser crea un nuovo utente
func CreateUser(user models.User) (*models.User, error) {
	log := utils.WithContext()
	// Registra un messaggio di log indicando che la creazione dell'utente è iniziata
	log.Info("Creo user --> request in ingresso:", user)

	// Chiama la funzione CreateUser del repository per inserire l'utente nel database
	// La funzione restituisce un risultato che contiene l'ID dell'utente appena creato e un eventuale errore
	result, err := repository.CreateUser(user)
	if err != nil {
		// Se si verifica un errore durante la creazione dell'utente, registra un messaggio di log e restituisce l'errore
		log.Errorf("Error durante la create user: %v", err)
		return nil, err
	}

	// Converte l'ID inserito (InsertedID) in una stringa
	// InsertedID è di tipo interface{}, quindi deve essere convertito a primitive.ObjectID
	userID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		// Se la conversione fallisce, restituisce un errore
		return nil, errors.New("errore durante la conversione dell InsertedID to ObjectID")
	}

	// Imposta l'ID dell'utente nella struttura user con il valore dell'ID convertito in stringa
	//MongoDB utilizza per identificare univocamente i documenti.
	//Tuttavia, ObjectID non è direttamente leggibile come stringa normale, quindi il metodo Hex() viene utilizzato per convertirlo in una stringa esadecimale leggibile.
	user.ID = userID.Hex()

	log.Info("User creato con ID: %s", user.ID)

	// Restituisce un puntatore alla struttura user appena creata e nil come errore
	// Questo evita di copiare l'intera struttura, permette modifiche successive, e utilizza nil per indicare l'assenza di errore
	return &user, nil
	//Restituire un puntatore (&user) evita di copiare l'intera struttura user, il che è più efficiente in termini di memoria e prestazioni.

	//Mutabilità:
	//Con un puntatore, il chiamante della funzione può modificare direttamente i campi della struttura user senza dover lavorare con una copia separata.
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (*models.User, error) {
	log := utils.WithContext()

	log.Info("Cerco utente Id: %s", id)
	user, err := repository.GetUserByID(id)
	if err != nil {
		log.Printf("Error retrieving user by ID: %s, error: %v", id, err)
	}
	return user, err
}

// DeleteUserByID deletes a user by ID
func DeleteUserByID(id string) error {
	log := utils.WithContext()

	log.Info("Cancello utente con Id: %s", id)
	err := repository.DeleteUserByID(id)
	if err != nil {
		log.Errorf("Error deleting user by ID: %s, error: %v", id, err)
	}
	return err
}

// UpdateUser updates a user by ID
// UpdateUser aggiorna un utente tramite ID
func UpdateUser(id string, user models.User) (*models.User, error) {
	log := utils.WithContext()

	// Registra un messaggio di log indicando che l'aggiornamento dell'utente è iniziato
	log.Info("Service: Update utente con ID: %s, e request in ingresso: %v", id, user)

	// Chiama la funzione UpdateUser del repository per aggiornare l'utente nel database
	// La funzione restituisce un eventuale errore
	err := repository.UpdateUser(id, user)
	if err != nil {
		// Se si verifica un errore durante l'aggiornamento dell'utente, registra un messaggio di log e restituisce l'errore
		log.Errorf("Error updating user by ID: %s, error: %v", id, err)
		return nil, err
	}

	// Chiama la funzione GetUserByID del repository per ottenere l'utente aggiornato dal database
	// La funzione restituisce un puntatore all'utente aggiornato e un eventuale errore
	updatedUser, err := repository.GetUserByID(id)
	if err != nil {
		// Se si verifica un errore durante il recupero dell'utente aggiornato, registra un messaggio di log
		log.Errorf("Error retrieving updated user by ID: %s, error: %v", id, err)
	}

	// Restituisce il puntatore all'utente aggiornato e un eventuale errore
	return updatedUser, err
}
