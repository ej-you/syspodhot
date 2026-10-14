package middleware

import (
	"errors"
	"fmt"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
	gojwt "github.com/golang-jwt/jwt/v5"

	"syspodhot/config"
	"syspodhot/internal/app/entity"
	"syspodhot/internal/app/service/server/errhandler"
	"syspodhot/internal/pkg/jwt"
)

const (
	_accessTokenCookie = "access" // key to store access token in cookies

	_tokenCtxKey      = "token"      // key for the JWT-token string in the fiber ctx
	_userClaimsCtxKey = "userClaims" // key for the user claims in the fiber ctx
)

// RestrictArea returns 404 code if token cookie is not specified.
func RestrictArea() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// parse cookies from request
		accessToken := ctx.Cookies(_accessTokenCookie)

		// if cookies is not specified
		if accessToken == "" {
			return errhandler.CustomErrorHandler(ctx, fiber.ErrNotFound)
		}
		return ctx.Next()
	}
}

// AlreadyLoggedIn redirects to index page if token cookie is specified.
func AlreadyLoggedIn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// parse cookie from request
		accessToken := ctx.Cookies(_accessTokenCookie)

		// if cookie is specified
		if accessToken != "" {
			return ctx.Redirect("/", http.StatusSeeOther)
		}
		return ctx.Next()
	}
}

// JWTAccess parses access token from request header to context and validates it.
func JWTAccess(cfg *config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// parse token string from access cookie
		token := ctx.Cookies(_accessTokenCookie)
		if token == "" {
			return handleTokenErr(ctx, errors.New("missing"))
		}
		// parse user claims from token
		userClaims, err := parseUserClaims(cfg.Server.Auth.AccessToken.Secret, token)
		if err != nil {
			return handleTokenErr(ctx, err)
		}
		// save token string and user claims to fiber context
		ctx.Locals(_tokenCtxKey, token)
		ctx.Locals(_userClaimsCtxKey, userClaims)
		return ctx.Next()
	}
}

// handleTokenErr handles a token error using custom error handler.
func handleTokenErr(ctx *fiber.Ctx, err error) error {
	return errhandler.CustomErrorHandler(ctx,
		fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("token: %v", err)),
	)
}

// parseTokenWithClaims returns parsed token user claims.
func parseUserClaims(tokenSecret []byte, token string) (*entity.UserClaims, error) {
	claims := &jwt.TokenClaims[*entity.UserClaims]{}
	// parse token object with claims from token string
	_, err := gojwt.ParseWithClaims(token, claims, func(t *gojwt.Token) (any, error) {
		if t.Method != gojwt.SigningMethodHS256 {
			return nil, errors.New("invalid signature")
		}
		return tokenSecret, nil
	})
	// if token is expired
	if errors.Is(err, gojwt.ErrTokenExpired) {
		return nil, errors.New("parse: token is expired")
	}
	// other error
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return claims.ExtraClaims, nil
}
