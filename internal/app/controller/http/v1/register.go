package v1

import "github.com/gofiber/fiber/v2"

// RegisterClientEndpoints registers all endpoints for client controller.
func RegisterClientEndpoints(router fiber.Router, ctrlClient *ControllerClient,
	alreadyLoggedIn fiber.Handler) {

	router.Get("/", ctrlClient.Index)
	router.Get("/album/:id", ctrlClient.Album)

	router.Get("/admin", alreadyLoggedIn, ctrlClient.Auth)
	router.Post("/admin", ctrlClient.AuthAPI)
}

// RegisterAdminEndpoints registers all endpoints for admin controller.
func RegisterAdminEndpoints(router fiber.Router, ctrlAdmin *ControllerAdmin,
	restrictArea, mwJWTAccess fiber.Handler) {

	admin := router.Group("/admin-area", restrictArea, mwJWTAccess)

	admin.Get("/album/new", ctrlAdmin.AlbumCreate)
}
