package db

import (
	"context"
	goerrors "errors"
	"fmt"

	"syspodhot/internal/app/entity"
	"syspodhot/internal/app/errors"
	"syspodhot/internal/pkg/db"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _collectionUser = (user{}).collectionName()

// UserRepoDB represents a user DB repo.
type UserRepoDB struct {
	dbStorage *db.DB
}

// NewUserRepoDB returns a new instance of UserRepoDB.
func NewUserRepoDB(dbStorage *db.DB) *UserRepoDB {
	return &UserRepoDB{
		dbStorage: dbStorage,
	}
}

// Create creates a new user record.
func (r *UserRepoDB) Create(ctx context.Context, user *entity.User) error {
	userDTO := newUser(user)
	// check that user with given username is not exists
	_, err := r.GetByUsername(ctx, userDTO.Username)
	if err == nil { // NOT err
		return fmt.Errorf("get by username: %w", errors.ErrAlreadyExists)
	}
	if !goerrors.Is(err, errors.ErrNotFound) {
		return err
	}

	_, err = r.dbStorage.
		Collection(_collectionUser).
		InsertOne(ctx, userDTO)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}

// GetByUsername returns user with given username.
func (r *UserRepoDB) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := &user{}

	err := r.dbStorage.
		Collection(_collectionUser).
		FindOne(ctx, bson.M{"username": username}).
		Decode(user)
	if goerrors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user.toEntity(), nil
}

// GetAll returns all users.
func (r *UserRepoDB) GetAll(ctx context.Context) ([]*entity.User, error) {
	userListRaw := make([]user, 0)

	cursor, err := r.dbStorage.
		Collection(_collectionUser).
		Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("get cursor: %w", err)
	}
	if err := cursor.All(ctx, &userListRaw); err != nil {
		return nil, fmt.Errorf("decode results: %w", err)
	}

	userList := make([]*entity.User, len(userListRaw))
	for idx := range userListRaw {
		userList[idx] = userListRaw[idx].toEntity()
	}
	return userList, nil
}

// DeleteByUsername deletes user with given username.
func (r *UserRepoDB) DeleteByUsername(ctx context.Context, username string) error {
	user := &user{}

	err := r.dbStorage.
		Collection(_collectionUser).
		FindOneAndDelete(ctx, bson.M{"username": username}).
		Decode(user)
	if goerrors.Is(err, mongo.ErrNoDocuments) {
		return errors.ErrNotFound
	}
	return err // err OR nil
}
