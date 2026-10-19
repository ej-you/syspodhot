package v1

import "github.com/gofiber/fiber/v2"

// RegisterClientEndpoints registers all endpoints for client controller.
func RegisterClientEndpoints(router fiber.Router, ctrlClient *ControllerClient,
	alreadyLoggedIn fiber.Handler) {

	router.Get("/", ctrlClient.Index)
	router.Get("/album/:id", ctrlClient.Album)
	router.Get("/concerts", ctrlClient.Concerts)

	router.Get("/login", alreadyLoggedIn, ctrlClient.Auth)
	router.Post("/login", ctrlClient.AuthAPI)
}

// RegisterAdminEndpoints registers all endpoints for admin controller.
func RegisterAdminEndpoints(router fiber.Router, ctrlAdmin *ControllerAdmin,
	restrictArea, mwJWTAccess fiber.Handler) {

	admin := router.Group("/admin", restrictArea, mwJWTAccess)

	admin.Get("/album/new", ctrlAdmin.AlbumCreate)
	admin.Get("/album/:id", ctrlAdmin.AlbumEdit)
	admin.Post("/album/new", ctrlAdmin.AlbumCreateAPI)
	admin.Post("/album/:id", ctrlAdmin.AlbumEditAPI)
	admin.Post("/album/:id/delete", ctrlAdmin.AlbumDeleteAPI)

	admin.Get("/concert/new", ctrlAdmin.ConcertCreate)
	admin.Get("/concert/:id", ctrlAdmin.ConcertEdit)
	admin.Post("/concert/new", ctrlAdmin.ConcertCreateAPI)
	admin.Post("/concert/:id", ctrlAdmin.ConcertEditAPI)
	admin.Post("/concert/:id/delete", ctrlAdmin.ConcertDeleteAPI)
}
