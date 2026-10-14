package v1

import (
	"github.com/gofiber/fiber/v2"

	"syspodhot/internal/app/usecase"
	"syspodhot/internal/pkg/validator"
)

// ControllerAdmin represents an http-controller for admin cases.
type ControllerAdmin struct {
	albumUC   *usecase.AlbumUC
	concertUC *usecase.ConcertUC
	valid     validator.Validator
}

// NewControllerAdmin returns a new instance of ControllerAdmin.
func NewControllerAdmin(albumUC *usecase.AlbumUC, concertUC *usecase.ConcertUC,
	valid validator.Validator) *ControllerAdmin {

	return &ControllerAdmin{
		albumUC:   albumUC,
		concertUC: concertUC,
		valid:     valid,
	}
}

// AlbumCreate renders HTML with new album page.
func (c *ControllerAdmin) AlbumCreate(ctx *fiber.Ctx) error {
	return ctx.SendString("temp answer")
}

// AlbumEdit renders HTML with existing album page.
func (c *ControllerAdmin) AlbumEdit(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// AlbumCreateAPI creates new album.
func (c *ControllerAdmin) AlbumCreateAPI(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// AlbumEditAPI saves album changes.
func (c *ControllerAdmin) AlbumEditAPI(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// AlbumDeleteAPI deletes album.
func (c *ControllerAdmin) AlbumDeleteAPI(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// Concert renders HTML with new concert page.
func (c *ControllerAdmin) ConcertCreate(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// Concert renders HTML with existing concert page.
func (c *ControllerAdmin) ConcertEdit(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// ConcertCreateAPI creates new concert.
func (c *ControllerAdmin) ConcertCreateAPI(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// ConcertEditAPI saves concert changes.
func (c *ControllerAdmin) ConcertEditAPI(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// ConcertDeleteAPI deletes concert.
func (c *ControllerAdmin) ConcertDeleteAPI(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}
