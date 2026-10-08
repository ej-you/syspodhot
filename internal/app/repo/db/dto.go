package db

import (
	"uuid"

	"syspodhot/internal/app/entity"
)

type user struct {
	ID       uuid.UUID `bson:"_id"`
	Username string    `bson:"username"`
	Password []byte    `bson:"password"`
}

func newUser(userInst *entity.User) *user {
	return &user{
		ID:       userInst.ID,
		Username: userInst.Username,
		Password: userInst.Password,
	}
}

func (u *user) toEntity() *entity.User {
	return &entity.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}

func (user) collectionName() string {
	return "users"
}
