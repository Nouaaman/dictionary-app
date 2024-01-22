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

	// check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// get the collection
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	return &Dictionary{client: client, collection: collection}, nil
}

func (d *Dictionary) Add(word, definition string) error {
	// check if already exists
	existingEntry, err := d.Get(word)
	if err == nil && existingEntry != nil {
		return ErrWordAlreadyExists
	}

	entry := Entry{Word: word, Definition: definition}

	// insert into db
	_, err = d.collection.InsertOne(context.Background(), entry)
	if err != nil {
		return err
	}

	return nil
}

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

func (d *Dictionary) Remove(word string) error {
	_, err := d.collection.DeleteOne(context.Background(), bson.M{"word": word})
	if err == mongo.ErrNoDocuments {
		return ErrWordNotFound
	} else if err != nil {
		return err
	}

	return nil
}

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
