package usecase

import (
	"context"
	"fmt"
	"uuid"

	"syspodhot/internal/app/entity"
	repodb "syspodhot/internal/app/repo/db"
	"syspodhot/internal/pkg/password"
)

// UserUC represents a user usecase.
type UserUC struct {
	userRepoDB *repodb.UserRepoDB
}

// NewUserUC returns a new instance of UserUC.
func NewUserUC(userRepoDB *repodb.UserRepoDB) *UserUC {
	return &UserUC{
		userRepoDB: userRepoDB,
	}
}

// Create creates a new user in the DB and returns them.
func (u *UserUC) Create(username string, rawPasswd []byte) (*entity.User, error) {
	// hash password
	hashPasswd, err := password.Encode(rawPasswd)
	if err != nil {
		return nil, fmt.Errorf("encode password: %w", err)
	}
	userInst := &entity.User{
		ID:       uuid.New(),
		Username: username,
		Password: hashPasswd,
	}

	if err := u.userRepoDB.Create(context.Background(), userInst); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return userInst, nil
}

// DeleteByUsername deletes user with given username.
func (u *UserUC) DeleteByUsername(username string) error {
	if err := u.userRepoDB.DeleteByUsername(context.Background(), username); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// GetAll returns all users.
func (u *UserUC) GetAll() ([]*entity.User, error) {
	userList, err := u.userRepoDB.GetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}
	return userList, nil
}
