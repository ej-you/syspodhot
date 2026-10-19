package v1

import (
	goerrors "errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"

	"syspodhot/config"
	"syspodhot/internal/app/entity"
	"syspodhot/internal/app/errors"
	"syspodhot/internal/app/usecase"
	"syspodhot/internal/pkg/validator"
)

// ControllerAdmin represents an http-controller for admin cases.
type ControllerAdmin struct {
	albumUC         *usecase.AlbumUC
	concertUC       *usecase.ConcertUC
	valid           validator.Validator
	mediaAlbum      string
	mediaAlbumURL   string
	mediaConcert    string
	mediaConcertURL string
}

// NewControllerAdmin returns a new instance of ControllerAdmin.
func NewControllerAdmin(cfg *config.Config, albumUC *usecase.AlbumUC,
	concertUC *usecase.ConcertUC, valid validator.Validator) *ControllerAdmin {

	return &ControllerAdmin{
		albumUC:         albumUC,
		concertUC:       concertUC,
		valid:           valid,
		mediaAlbum:      cfg.Media.Album,
		mediaAlbumURL:   cfg.Media.AlbumURL,
		mediaConcert:    cfg.Media.Concert,
		mediaConcertURL: cfg.Media.ConcertURL,
	}
}

// AlbumCreate renders HTML with new album page.
func (c *ControllerAdmin) AlbumCreate(ctx *fiber.Ctx) error {
	return ctx.Render("admin_album", fiber.Map{
		"mode": "new",
	})
}

// AlbumEdit renders HTML with existing album page.
func (c *ControllerAdmin) AlbumEdit(ctx *fiber.Ctx) error {
	inputPath := &idPath{}
	if err := inputPath.Parse(ctx, c.valid); err != nil {
		return err
	}

	album, err := c.albumUC.Get(inputPath.ID)
	if goerrors.Is(err, errors.ErrNotFound) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("get album by id: %w", err)
	}

	return ctx.Render("admin_album", fiber.Map{
		"mode":       "edit",
		"album":      album,
		"mediaAlbum": c.mediaAlbumURL,
	})
}

// AlbumCreateAPI creates new album.
func (c *ControllerAdmin) AlbumCreateAPI(ctx *fiber.Ctx) error {
	inputBody := &albumBody{}
	if err := inputBody.Parse(ctx, c.valid, c.mediaAlbum); err != nil {
		return err
	}

	album := &entity.Album{
		Photo:      inputBody.Photo,
		Title:      inputBody.Title,
		Type:       inputBody.Type,
		Annotation: inputBody.Annotation,
		ReleasedAt: inputBody.ReleasedAt,
		Streaming: entity.Streaming{
			VKLink:     inputBody.VKLink,
			YandexLink: inputBody.YandexLink,
		},
		Songs: make([]entity.Song, len(inputBody.Songs)),
	}
	for idx := range inputBody.Songs {
		album.Songs[idx].Title = inputBody.Songs[idx].Title
		album.Songs[idx].Duration = inputBody.Songs[idx].Duration
	}

	if err := c.albumUC.Create(album); err != nil {
		return fmt.Errorf("create album: %w", err)
	}
	url := "/album/" + album.ID.String()
	return ctx.Redirect(url, fiber.StatusSeeOther)
}

// AlbumEditAPI saves album changes.
func (c *ControllerAdmin) AlbumEditAPI(ctx *fiber.Ctx) error {
	inputPath := &idPath{}
	if err := inputPath.Parse(ctx, c.valid); err != nil {
		return err
	}
	inputBody := &albumBody{}
	if err := inputBody.Parse(ctx, c.valid, c.mediaAlbum); err != nil {
		return err
	}

	albumOld, err := c.albumUC.Get(inputPath.ID)
	if err != nil {
		return fmt.Errorf("get old album data: %w", err)
	}
	// remove old photo
	if albumOld.Photo != "" && inputBody.Photo != "" {
		if err := os.Remove(path.Join(c.mediaAlbum, albumOld.Photo)); err != nil {
			slog.Error("album edit api", "remove old cover", err)
		}
	}

	album := &entity.Album{
		ID:         inputPath.ID,
		Photo:      inputBody.Photo,
		Title:      inputBody.Title,
		Type:       inputBody.Type,
		Annotation: inputBody.Annotation,
		ReleasedAt: inputBody.ReleasedAt,
		Streaming: entity.Streaming{
			VKLink:     inputBody.VKLink,
			YandexLink: inputBody.YandexLink,
		},
		Songs: make([]entity.Song, len(inputBody.Songs)),
	}
	for idx := range inputBody.Songs {
		album.Songs[idx].Title = inputBody.Songs[idx].Title
		album.Songs[idx].Duration = inputBody.Songs[idx].Duration
	}

	if err := c.albumUC.Update(album); err != nil {
		return fmt.Errorf("update album: %w", err)
	}
	url := "/album/" + album.ID.String()
	return ctx.Redirect(url, fiber.StatusSeeOther)
}

// AlbumDeleteAPI deletes album.
func (c *ControllerAdmin) AlbumDeleteAPI(ctx *fiber.Ctx) error {
	inputPath := &idPath{}
	if err := inputPath.Parse(ctx, c.valid); err != nil {
		return err
	}

	// remove record from DB
	album, err := c.albumUC.Delete(inputPath.ID)
	if err != nil {
		return fmt.Errorf("delete album: %w", err)
	}
	// remove photo
	if album.Photo != "" {
		if err := os.Remove(path.Join(c.mediaAlbum, album.Photo)); err != nil {
			slog.Error("album edit api", "remove album cover", err)
		}
	}

	return ctx.Redirect("/#music", fiber.StatusSeeOther)
}

// Concert renders HTML with new concert page.
func (c *ControllerAdmin) ConcertCreate(ctx *fiber.Ctx) error {
	return ctx.Render("admin_concert", fiber.Map{
		"mode": "new",
	})
}

// Concert renders HTML with existing concert page.
func (c *ControllerAdmin) ConcertEdit(ctx *fiber.Ctx) error {
	inputPath := &idPath{}
	if err := inputPath.Parse(ctx, c.valid); err != nil {
		return err
	}

	concert, err := c.concertUC.Get(inputPath.ID)
	if goerrors.Is(err, errors.ErrNotFound) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("get concert by id: %w", err)
	}

	return ctx.Render("admin_concert", fiber.Map{
		"mode":    "edit",
		"concert": concert,
	})
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
