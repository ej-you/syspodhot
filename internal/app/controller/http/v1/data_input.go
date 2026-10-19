package v1

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"
	"uuid"

	"github.com/gofiber/fiber/v2"

	utilsfile "syspodhot/internal/pkg/utils/file"
	"syspodhot/internal/pkg/validator"
)

const (
	_fmtDate = "2006-01-02"
	_minute  = 60
)

var _patternMMSS = regexp.MustCompile(`\d{1,3}:[0-5]\d`)

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

// idPath represents a data with album/concert ID in path params.
type idPath struct {
	// album/concert ID
	ID uuid.UUID `params:"id" validate:"required"`
}

// Parse parses idPath request data and validates it.
func (i *idPath) Parse(ctx *fiber.Ctx, valid validator.Validator) error {
	// parse path-params
	if err := ctx.ParamsParser(i); err != nil {
		return err
	}
	// validate parsed data
	if err := valid.Validate(i); err != nil {
		return err
	}
	return nil
}

// albumBody represents an album data.
type albumBody struct {
	Title         string     `form:"title" validate:"required"`
	Type          string     `form:"type" validate:"required,oneof=Альбом Сингл EP"`
	Annotation    string     `form:"annotation" validate:"omitempty"`
	ReleasedAtRaw string     `form:"released_at" validate:"required"`
	VKLink        string     `form:"streaming__vk" validate:"omitempty"`
	YandexLink    string     `form:"streaming__yandex" validate:"omitempty"`
	Songs         []songBody `form:"songs" validate:"dive"`

	Photo      string    `form:"-" validate:"-"`
	ReleasedAt time.Time `form:"-" validate:"-"`
}

// songBody represents a song data for album.
type songBody struct {
	Title       string `form:"title" validate:"required"`
	DurationRaw string `form:"duration" validate:"required"`

	Duration int `form:"-" validate:"-"`
}

// Parse parses idPath request data and validates it.
func (a *albumBody) Parse(ctx *fiber.Ctx, valid validator.Validator, mediaDir string) error {
	var err error

	if err := ctx.BodyParser(a); err != nil {
		return fmt.Errorf("parse body: %w", err)
	}
	// parse release date
	a.ReleasedAt, err = time.Parse(_fmtDate, a.ReleasedAtRaw)
	if err != nil {
		return fmt.Errorf("%w: parse release date: %v", validator.ErrValidate, err)
	}
	// parse song durations
	for idx := range a.Songs {
		a.Songs[idx].Duration, err = mmssToSeconds(a.Songs[idx].DurationRaw)
		if err != nil {
			return fmt.Errorf("%w: parse song duration: %v", validator.ErrValidate, err)
		}
	}
	// validate parsed data
	if err := valid.Validate(a); err != nil {
		return err
	}

	form, _ := ctx.MultipartForm()
	slog.Debug("parse album data", "form", form)

	// upload album cover
	a.Photo, err = utilsfile.UploadAndSaveFile(ctx, mediaDir, "cover", nil)
	if err != nil {
		return fmt.Errorf("upload and save photo: %w", err)
	}
	return nil
}

// mmssToSeconds parses mm:ss and converts to seconds.
func mmssToSeconds(raw string) (int, error) {
	if !_patternMMSS.MatchString(raw) {
		return 0, errors.New("expect mm:ss format")
	}
	sep := strings.Split(raw, ":")

	mins, err := strconv.Atoi(sep[0])
	if err != nil {
		return 0, fmt.Errorf("parse minutes: %w", err)
	}
	secs, err := strconv.Atoi(sep[1])
	if err != nil {
		return 0, fmt.Errorf("parse seconds: %w", err)
	}
	return mins*_minute + secs, nil
}
