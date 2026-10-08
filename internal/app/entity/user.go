package entity

import "uuid"

// User represents an app user with auth user data.
type User struct {
	// user id
	ID uuid.UUID `json:"id"`
	// user username
	Username string `json:"username"`
	// user password hash
	Password []byte `json:"-"`
}

// UserClaims represents a claims with user data for JWT-tokens.
type UserClaims struct {
	// user ID
	ID uuid.UUID `json:"id"`
	// user username
	Username string `json:"username"`
}

// UserWithToken is a user object and a token string.
type UserWithToken struct {
	User        *User
	AccessToken string `json:"-"`
}
