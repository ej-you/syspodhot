package v1

import (
	"uuid"

	"github.com/gofiber/fiber/v2"

	"syspodhot/internal/pkg/validator"
)

// authBody represents a data for user auth (log in).
type authBody struct {
	// user username
	Username string `json:"username" validate:"required,max=50"`
	// user password
	Password string `json:"password" validate:"required,min=8,max=40"`
}

// Parse parses authBody request data and validates it.
func (a *authBody) Parse(ctx *fiber.Ctx, valid validator.Validator) error {
	// parse body
	if err := ctx.BodyParser(a); err != nil {
		return err
	}
	// validate parsed data
	if err := valid.Validate(a); err != nil {
		return err
	}
	return nil
}

// albumPath represents a data with album ID in path params.
type albumPath struct {
	// album ID
	ID uuid.UUID `params:"id" validate:"required"`
}

// Parse parses albumPath request data and validates it.
func (a *albumPath) Parse(ctx *fiber.Ctx, valid validator.Validator) error {
	// parse path-params
	if err := ctx.ParamsParser(a); err != nil {
		return err
	}
	// validate parsed data
	if err := valid.Validate(a); err != nil {
		return err
	}
	return nil
}
