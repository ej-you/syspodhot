package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	_closeTimeout = 10 * time.Second // timeout to close DB connection
	_pingTimeout  = 5 * time.Second  // ping timeout to check DB is ready
)

// DB represents a DB connection.
type DB struct {
	*mongo.Database
}

// New returns a new instance of DB with connection to the given DSN and database.
func New(dsn, dbName string) (*DB, error) {
	opts := options.Client().ApplyURI(dsn)
	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("conn to db: %w", err)
	}

	db := &DB{mongoClient.Database(dbName)}
	// check connection
	if err := db.IsReady(); err != nil {
		return nil, err
	}
	return db, nil
}

// IsReady checks that DB is ready to use.
// Returns nil error if DB is ready.
func (db *DB) IsReady() error {
	ctx, cancel := context.WithTimeout(context.Background(), _pingTimeout)
	defer cancel()

	if err := db.Client().Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("db is unreachable: %w", err)
	}
	return nil
}

// Close closes db connection.
func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), _closeTimeout)
	defer cancel()

	return db.Client().Disconnect(ctx)
}
