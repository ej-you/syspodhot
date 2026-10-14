// Package errhandler contains custom error handler func for fiber server.
package errhandler

import (
	goerrors "errors"
	"log/slog"

	fiber "github.com/gofiber/fiber/v2"

	"syspodhot/internal/pkg/validator"
)

// CustomErrorHandler represents a custom error handler for server.
func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	var (
		errCode int
		errMsg  string
		errLog  bool
	)

	var fiberErr *fiber.Error

	switch {
	// *fiber.Error error
	case goerrors.As(err, &fiberErr):
		errCode = fiberErr.Code
		if fiberErr.Code == fiber.StatusNotFound {
			errMsg = "страница не найдена"
			errLog = false
		} else {
			errMsg = fiberErr.Message
			errLog = true
		}
	// validate error
	case goerrors.Is(err, validator.ErrValidate):
		errCode = fiber.StatusBadRequest
		errMsg = err.Error()
		errLog = false
	// unknown error
	default:
		errCode = fiber.StatusInternalServerError
		errMsg = "внутренняя ошибка сервера"
		errLog = true
	}

	if errLog {
		slog.Error("error handler",
			"code", errCode,
			"cause", err)
	}

	return ctx.Status(errCode).Render("error", fiber.Map{
		"code":    errCode,
		"message": errMsg,
	})
}
