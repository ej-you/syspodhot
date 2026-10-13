package v1

import "github.com/gofiber/fiber/v2"

// RegisterClientEndpoints registers all endpoints for client controller.
func RegisterClientEndpoints(router fiber.Router, ctrlClient *ControllerClient) {
	router.Get("/", ctrlClient.Index)
	router.Get("/album/:id", ctrlClient.Album)
}
