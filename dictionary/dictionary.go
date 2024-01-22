// dictionary/dictionary.go

package dictionary

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrWordAlreadyExists = errors.New("word already exists")
	ErrWordNotFound      = errors.New("word not found")
)

type Entry struct {
	Word       string `json:"word" bson:"word"`
	Definition string `json:"definition" bson:"definition"`
}

type Dictionary struct {
	client     *mongo.Client
	collection *mongo.Collection
}

const (
	databaseName   = "dictionary"
	collectionName = "words"
)

// New initializes a new MongoDB-backed dictionary
func New() (*Dictionary, error) {

	connectionString := "mongodb+srv://nouaaman:v87QnLnYoTv6byLQ@dictionary.cemt71j.mongodb.net/?retryWrites=true&w=majority"

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Get the collection
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	return &Dictionary{client: client, collection: collection}, nil
}

// Add adds a new entry to the dictionary
func (d *Dictionary) Add(word, definition string) error {
	// Check if the word already exists
	existingEntry, err := d.Get(word)
	if err == nil && existingEntry != nil {
		return ErrWordAlreadyExists
	}

	// Create a new entry
	entry := Entry{Word: word, Definition: definition}

	// Insert the entry into MongoDB
	_, err = d.collection.InsertOne(context.Background(), entry)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an entry from the dictionary by word
func (d *Dictionary) Get(word string) (*Entry, error) {
	var entry Entry
	err := d.collection.FindOne(context.Background(), bson.M{"word": word}).Decode(&entry)

	if err == mongo.ErrNoDocuments {
		return nil, ErrWordNotFound
	} else if err != nil {
		return nil, err
	}

	return &entry, nil
}

// Remove removes an entry from the dictionary by word
func (d *Dictionary) Remove(word string) error {
	_, err := d.collection.DeleteOne(context.Background(), bson.M{"word": word})
	if err == mongo.ErrNoDocuments {
		return ErrWordNotFound
	} else if err != nil {
		return err
	}

	return nil
}

// List retrieves all entries from the dictionary
func (d *Dictionary) List() ([]Entry, error) {
	var entries []Entry
	cursor, err := d.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
