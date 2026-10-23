package v1

import (
	goerrors "errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"slices"
	"time"
	"uuid"

	"github.com/gofiber/fiber/v2"

	"syspodhot/config"
	"syspodhot/internal/app/entity"
	"syspodhot/internal/app/errors"
	"syspodhot/internal/app/usecase"
	"syspodhot/internal/pkg/validator"
)

const _dirMode = 0o750

// ControllerAdmin represents an http-controller for admin cases.
type ControllerAdmin struct {
	albumUC         *usecase.AlbumUC
	concertUC       *usecase.ConcertUC
	valid           validator.Validator
	locale          *time.Location
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
		locale:          cfg.Locale,
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
		defer os.Remove(path.Join(c.mediaAlbum, album.Photo)) // remove uploaded file
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
			slog.Error("album delete api", "remove album cover", err)
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
		"mode":         "edit",
		"concert":      concert,
		"mediaConcert": c.mediaConcertURL,
	})
}

// ConcertCreateAPI creates new concert.
func (c *ControllerAdmin) ConcertCreateAPI(ctx *fiber.Ctx) error {
	photosFolder := uuid.New().String()
	filesDir := path.Join(c.mediaConcert, photosFolder)
	// create concert subdir
	err := os.Mkdir(filesDir, _dirMode)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("create photos subdir: %w", err)
	}

	inputBody := &concertBody{}
	if err := inputBody.Parse(ctx, c.valid, filesDir, c.locale); err != nil {
		defer os.RemoveAll(filesDir) // remove uploaded files
		return err
	}

	concert := &entity.Concert{
		Title:        inputBody.Title,
		Place:        inputBody.Place,
		Datetime:     inputBody.Datetime,
		TicketLink:   inputBody.TicketLink,
		PhotosFolder: photosFolder,
		Photos:       inputBody.Photos,
		Video:        inputBody.Video,
	}

	if err := c.concertUC.Create(concert); err != nil {
		defer os.RemoveAll(filesDir) // remove uploaded files
		return fmt.Errorf("create concert: %w", err)
	}
	return ctx.Redirect("/concerts", fiber.StatusSeeOther)
}

// ConcertEditAPI saves concert changes.
func (c *ControllerAdmin) ConcertEditAPI(ctx *fiber.Ctx) error {
	inputPath := &idPath{}
	if err := inputPath.Parse(ctx, c.valid); err != nil {
		return err
	}
	concertOld, err := c.concertUC.Get(inputPath.ID)
	if err != nil {
		return fmt.Errorf("get old concert data: %w", err)
	}

	filesDir := path.Join(c.mediaConcert, concertOld.PhotosFolder)
	inputBody := &concertBody{}
	if err := inputBody.Parse(ctx, c.valid, filesDir, c.locale); err != nil {
		return err
	}

	// remove selected photos
	for _, photo := range inputBody.DeletePhotos {
		if err := os.Remove(path.Join(filesDir, photo)); err != nil {
			slog.Error("concert edit api", "remove photo", err)
		}
		// delete filename from slice
		idx := slices.Index(concertOld.Photos, photo)
		if idx != -1 {
			concertOld.Photos = slices.Delete(concertOld.Photos, idx, idx+1)
		}
	}

	concert := &entity.Concert{
		ID:           inputPath.ID,
		Title:        inputBody.Title,
		Place:        inputBody.Place,
		Datetime:     inputBody.Datetime,
		TicketLink:   inputBody.TicketLink,
		PhotosFolder: concertOld.PhotosFolder,
		Photos:       append(concertOld.Photos, inputBody.Photos...),
		Video:        inputBody.Video,
	}

	if err := c.concertUC.Update(concert); err != nil {
		return fmt.Errorf("update concert: %w", err)
	}
	return ctx.Redirect("/concerts", fiber.StatusSeeOther)
}

// ConcertDeleteAPI deletes concert.
func (c *ControllerAdmin) ConcertDeleteAPI(ctx *fiber.Ctx) error {
	inputPath := &idPath{}
	if err := inputPath.Parse(ctx, c.valid); err != nil {
		return err
	}

	// remove record from DB
	concert, err := c.concertUC.Delete(inputPath.ID)
	if err != nil {
		return fmt.Errorf("delete concert: %w", err)
	}
	// remove photos folder
	if len(concert.Photos) != 0 {
		if err := os.RemoveAll(path.Join(c.mediaConcert, concert.PhotosFolder)); err != nil {
			slog.Error("concert delete api", "remove concert photos", err)
		}
	}
	return ctx.Redirect("/concerts", fiber.StatusSeeOther)
}
