package v1

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"syspodhot/internal/app/usecase"
	"syspodhot/internal/pkg/validator"
)

// ControllerClient represents an http-controller for client cases.
type ControllerClient struct {
	userUC    *usecase.UserUC
	albumUC   *usecase.AlbumUC
	concertUC *usecase.ConcertUC
	valid     validator.Validator
}

// NewControllerClient returns a new instance of ControllerClient.
func NewControllerClient(userUC *usecase.UserUC, albumUC *usecase.AlbumUC,
	concertUC *usecase.ConcertUC, valid validator.Validator) *ControllerClient {

	return &ControllerClient{
		userUC:    userUC,
		albumUC:   albumUC,
		concertUC: concertUC,
		valid:     valid,
	}
}

// Index renders HTML with main page.
func (c *ControllerClient) Index(ctx *fiber.Ctx) error {
	albums, err := c.albumUC.GetAll()
	if err != nil {
		return fmt.Errorf("get all albums: %w", err)
	}

	concerts, err := c.concertUC.GetComming()
	if err != nil {
		return fmt.Errorf("get comming concerts: %w", err)
	}

	return ctx.Render("index", fiber.Map{
		"albums":   albums,
		"concerts": concerts,
	})
}

// Album renders HTML with album page.
func (c *ControllerClient) Album(ctx *fiber.Ctx) error {
	inputPath := &albumPath{}
	if err := inputPath.Parse(ctx, c.valid); err != nil {
		return err
	}

	album, err := c.albumUC.Get(inputPath.ID)
	if err != nil {
		return fmt.Errorf("get album by id: %w", err)
	}

	return ctx.Render("album", fiber.Map{
		"album": album,
	})
}

// Concerts renders HTML with concerts page.
func (c *ControllerClient) Concerts(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// Auth renders HTML with auth page for admin.
func (c *ControllerClient) Auth(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}

// AuthAPI auths admin user.
func (c *ControllerClient) AuthAPI(ctx *fiber.Ctx) error {
	// TODO: implement
	panic("unimplemented")
}
