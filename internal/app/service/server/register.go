package server

import (
	"syspodhot/config"
	httpv1 "syspodhot/internal/app/controller/http/v1"
	repodb "syspodhot/internal/app/repo/db"
	"syspodhot/internal/app/usecase"
	"syspodhot/internal/pkg/db"
	"syspodhot/internal/pkg/validator"
)

// registerEndpointsV1 register all endpoints for 1st version of API.
func (s *Server) registerEndpointsV1(cfg *config.Config, dbStorage *db.DB,
	valid validator.Validator) {

	// set up static
	s.fiberApp.Static("/static", "./web/static")

	// create repos
	userRepoDB := repodb.NewUserRepoDB(dbStorage)
	albumRepoDB := repodb.NewAlbumRepoDB(dbStorage)
	concertRepoDB := repodb.NewConcertRepoDB(dbStorage)
	// create usecases
	userUC := usecase.NewUserUC(userRepoDB)
	albumUC := usecase.NewAlbumUC(albumRepoDB)
	concertUC := usecase.NewConcertUC(concertRepoDB)
	// create controllers
	controllerClient := httpv1.NewControllerClient(userUC, albumUC, concertUC, valid)

	// // middlewares
	// loginRequired := middleware.LoginRequired()
	// alreadyLoggedIn := middleware.AlreadyLoggedIn()
	// mwJWTAccess := middleware.JWTAccess(cfg)

	// register endpoints
	httpv1.RegisterClientEndpoints(s.fiberApp, controllerClient)
}
