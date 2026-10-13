package usecase

import (
	"context"
	goerrors "errors"
	"fmt"

	"syspodhot/config"
	"syspodhot/internal/app/entity"
	"syspodhot/internal/app/errors"
	repodb "syspodhot/internal/app/repo/db"
	"syspodhot/internal/pkg/jwt"
	"syspodhot/internal/pkg/password"
)

// AuthUC represents a user auth usecase.
type AuthUC struct {
	userRepoDB *repodb.UserRepoDB
	jwtBuilder jwt.Builder
}

// NewAuthUC returns a new instance of AuthUC.
func NewAuthUC(cfg *config.Config, userRepoDB *repodb.UserRepoDB) *AuthUC {
	return &AuthUC{
		userRepoDB: userRepoDB,
		jwtBuilder: *jwt.NewBuilder(
			cfg.Auth.AccessToken.Secret,
			nil,
			jwt.WithAccessTTL(cfg.Auth.AccessToken.TTL)),
	}
}

// LogIn authenticates user and obtains auth token.
func (u *AuthUC) LogIn(username string, rawPasswd []byte) (*entity.UserWithToken, error) {
	// get user object from DB with username
	userInst, err := u.userRepoDB.GetByUsername(context.Background(), username)
	if goerrors.Is(err, errors.ErrNotFound) {
		return nil, errors.ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	// check entered password is correct
	if !password.IsCorrect(rawPasswd, userInst.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	userClaims := &entity.UserClaims{
		ID:       userInst.ID,
		Username: userInst.Username,
	}
	// generate token
	token, err := u.jwtBuilder.ObtainAccess(userClaims)
	if err != nil {
		return nil, fmt.Errorf("obtain token: %w", err)
	}
	return &entity.UserWithToken{
		User:        userInst,
		AccessToken: token,
	}, nil
}
