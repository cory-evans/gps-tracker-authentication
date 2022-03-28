package databasev1

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabaseClient(MONGO_URI string, ctx context.Context) (*mongo.Client, error) {
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))

	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	return db, nil
}
