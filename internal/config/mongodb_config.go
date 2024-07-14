package config

import (
	"context"
	"myapp/internal/utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mongoClientInstance *mongo.Client
	databaseInstance    *mongo.Database
	once                sync.Once
)

// GetMongoClient ritorna l'istanza singleton del client MongoDB
func GetMongoClient() *mongo.Client {
	once.Do(func() {
		loadConfig()
	})
	return mongoClientInstance
}

// GetDatabase ritorna l'istanza singleton del database MongoDB
func GetDatabase() *mongo.Database {
	once.Do(func() {
		loadConfig()
	})
	return databaseInstance
}

// loadConfig carica la configurazione e stabilisce la connessione con MongoDB
func loadConfig() {
	log := utils.WithContext().WithField("function", "loadConfig")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := utils.EnvOrDefault("MONGO_URI", "mongodb://localhost:27017")

	var err error
	mongoClientInstance, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClientInstance.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Error on establishing a connection with the cluster: %v", err)
	}

	databaseName := utils.EnvOrDefault("MONGO_DATABASE", "myapp")
	databaseInstance = mongoClientInstance.Database(databaseName)
	log.Infof("Connected to MongoDB at %s", mongoURI)
}
