package repository

import (
	"context"
	"myapp/internal/config"
	"myapp/internal/models"
	"myapp/internal/utils"
	"myapp/internal/utils/constants"

	"github.com/openzipkin/zipkin-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUsers retrieves all users from the MongoDB collection
func GetUsers(ctx context.Context, tracer *zipkin.Tracer) ([]models.User, error) {
	log := utils.WithContext().WithField("function", "GetUsers_Repo")

	// Recupera lo span dal contesto
	parentSpan := zipkin.SpanFromContext(ctx)
	// Crea uno span figlio
	childSpan := tracer.StartSpan("GetUsers_Repo", zipkin.Parent(parentSpan.Context()))
	defer childSpan.Finish()

	// Crea un nuovo contesto con lo span figlio
	childCtx := zipkin.NewContext(ctx, childSpan)

	var users []models.User
	cursor, err := config.GetDatabase().Collection(constants.USERSCOLLECTION).Find(childCtx, bson.M{})
	if err != nil {
		log.Errorf("Error finding users: %v", err)
		return nil, err
	}

	if err = cursor.All(childCtx, &users); err != nil {
		log.Errorf("Error decoding users: %v", err)
		return nil, err
	}
	return users, nil
}

// CreateUser inserisce un nuovo utente nella collezione MongoDB
func CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	log := utils.WithContext().WithField("function", "CreateUser")
	db := config.GetDatabase()

	// Esegue l'operazione di inserimento e restituisce il risultato dell'inserimento e un eventuale errore
	result, err := db.Collection(constants.USERSCOLLECTION).InsertOne(context.TODO(), user)
	if err != nil {
		log.Errorf("Error creating user: %v", err)
	}
	return result, err
}

// GetUserByID recupera un utente per ID dalla collezione MongoDB
func GetUserByID(id string) (*models.User, error) {
	log := utils.WithContext().WithField("function", "GetUserByID")
	db := config.GetDatabase()

	// Converte l'ID esadecimale (stringa) in un ObjectID di MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Error converting ID: %v", err)
		return nil, err
	}

	// Crea una variabile per memorizzare l'utente trovato
	var user models.User

	// Esegue la query per trovare il documento con l'ObjectID specificato
	err = db.Collection(constants.USERSCOLLECTION).FindOne(context.TODO(), bson.M{constants.DOCUMENT_ID: objectID}).Decode(&user)
	if err != nil {
		log.Errorf("Error finding user by ID: %v", err)
		return nil, err
	}

	// Restituisce un puntatore all'utente trovato e nil come errore
	return &user, nil
}

// DeleteUserByID deletes a user by ID from the MongoDB collection
func DeleteUserByID(id string) error {
	log := utils.WithContext().WithField("function", "DeleteUserByID")
	db := config.GetDatabase()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Error converting ID: %v", err)
		return err
	}
	_, err = db.Collection(constants.USERSCOLLECTION).DeleteOne(context.TODO(), bson.M{constants.DOCUMENT_ID: objectID})
	if err != nil {
		log.Errorf("Error deleting user by ID: %v", err)
	}
	return err
}

// UpdateUser updates a user by ID in the MongoDB collection
func UpdateUser(id string, user models.User) error {
	log := utils.WithContext().WithField("function", "UpdateUser")
	db := config.GetDatabase()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Error converting ID: %v", err)
		return err
	}
	_, err = db.Collection(constants.USERSCOLLECTION).UpdateOne(
		context.TODO(),
		bson.M{constants.DOCUMENT_ID: objectID},
		bson.M{constants.SET: user},
	)
	if err != nil {
		log.Errorf("Error updating user by ID: %v", err)
	}
	return err
}
