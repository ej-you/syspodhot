package v1

import (
	goerrors "errors"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"

	"syspodhot/config"
	"syspodhot/internal/app/errors"
	"syspodhot/internal/app/usecase"
	"syspodhot/internal/pkg/cookie"
	"syspodhot/internal/pkg/validator"
)

const _accessTokenCookie = "access" // key to store access token in cookies

// ControllerClient represents an http-controller for client cases.
type ControllerClient struct {
	authUC              *usecase.AuthUC
	albumUC             *usecase.AlbumUC
	concertUC           *usecase.ConcertUC
	valid               validator.Validator
	accessCookieBuilder *cookie.Builder
}

// NewControllerClient returns a new instance of ControllerClient.
func NewControllerClient(cfg *config.Config, authUC *usecase.AuthUC, albumUC *usecase.AlbumUC,
	concertUC *usecase.ConcertUC, valid validator.Validator) *ControllerClient {

	return &ControllerClient{
		authUC:    authUC,
		albumUC:   albumUC,
		concertUC: concertUC,
		valid:     valid,
		accessCookieBuilder: cookie.NewBuilder(cfg.Auth.AccessToken.TTL,
			cookie.WithPath(cfg.Auth.AccessToken.Cookie.Path),
			cookie.WithSecure(cfg.Auth.AccessToken.Cookie.Secure),
			cookie.WithHTTPOnly(cfg.Auth.AccessToken.Cookie.HTTPOnly),
			cookie.WithSameSite(cfg.Auth.AccessToken.Cookie.SameSite),
		),
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
		"isAdmin":  ctx.Cookies(_accessTokenCookie) != "",
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
	if goerrors.Is(err, errors.ErrNotFound) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("get album by id: %w", err)
	}

	return ctx.Render("album", fiber.Map{
		"isAdmin": ctx.Cookies(_accessTokenCookie) != "",
		"album":   album,
	})
}

// Concerts renders HTML with concerts page.
func (c *ControllerClient) Concerts(ctx *fiber.Ctx) error {
	comming, err := c.concertUC.GetComming()
	if err != nil {
		return fmt.Errorf("get comming concerts: %w", err)
	}
	past, err := c.concertUC.GetPastGroupedByYear()
	if err != nil {
		return fmt.Errorf("get past concerts: %w", err)
	}

	return ctx.Render("concerts", fiber.Map{
		"isAdmin": ctx.Cookies(_accessTokenCookie) != "",
		"comming": comming,
		"past":    past,
	})
}

// Auth renders HTML with auth page for admin.
func (c *ControllerClient) Auth(ctx *fiber.Ctx) error {
	return ctx.Render("login", fiber.Map{})
}

// AuthAPI auths admin user.
func (c *ControllerClient) AuthAPI(ctx *fiber.Ctx) error {
	inputBody := &authBody{}
	if err := inputBody.Parse(ctx, c.valid); err != nil {
		return err
	}

	// log in existing user
	userWithToken, err := c.authUC.LogIn(inputBody.Username, []byte(inputBody.Password))
	if goerrors.Is(err, errors.ErrInvalidCredentials) {
		return fiber.NewError(400, "неверный логин или пароль")
	}
	if err != nil {
		return err
	}
	// add access token to cookies
	ctx.Cookie(c.accessCookieBuilder.Create(_accessTokenCookie, userWithToken.AccessToken))

	redirectURL, err := url.QueryUnescape(ctx.Query("next", "/"))
	if err != nil {
		return fmt.Errorf("unescape next param: %w", err)
	}
	return ctx.Redirect(redirectURL, fiber.StatusSeeOther)
}
