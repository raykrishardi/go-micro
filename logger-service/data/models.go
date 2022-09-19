package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string    `bson:"name" json:"name"`
	Data       string    `bson:"data" json:"data"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdateddAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	// If the collection doesn't exist then it will actually create it
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:       entry.Name,
		Data:       entry.Data,
		CreatedAt:  time.Now(),
		UpdateddAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting to logs:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	// Use bson.D (slice) if the order matters because we are sorting by created_at
	// Use bson.M (map) if the order does NOT matter
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	logs := []*LogEntry{}

	for cursor.Next(ctx) {
		item := LogEntry{}
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Log decoding log into slice:", err)
			return nil, err
		}
		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Convert the string ID to mongo ID
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	entry := LogEntry{}

	// Use bson.D (slice) if the order matters because we are sorting by created_at
	// Use bson.M (map) if the order does NOT matter
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Dropping the collection, next time, you call this, it will create the collection
	collection := client.Database("logs").Collection("logs")
	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Convert the string ID to mongo ID
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
