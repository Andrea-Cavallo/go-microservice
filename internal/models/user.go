package models

// User Combinazione di json e bson
// Utilizzando entrambe le annotazioni, puoi garantire che la stessa struttura User possa
// essere utilizzata senza problemi sia per la comunicazione API in formato JSON che per la memorizzazione e il recupero dei dati in MongoDB in formato BSON.
type User struct {
	ID    string `json:"id" bson:"_id,omitempty"` // _id,omitempty" specifica che il campo ID è mappato al campo _id in MongoDB e che deve essere omesso se vuoto (omitempty).
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}
