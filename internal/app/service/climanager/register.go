package climanager

import (
	"syspodhot/internal/app/controller/cli"
	repodb "syspodhot/internal/app/repo/db"
	"syspodhot/internal/app/usecase"
	"syspodhot/internal/pkg/db"
)

// registerCommands register all cmd manager commands.
func (c *CLIManager) registerCommands(dbStorage *db.DB) {
	// create repos
	userRepoDB := repodb.NewUserRepoDB(dbStorage)
	// create usecases
	userUC := usecase.NewUserUC(userRepoDB)
	// create controllers
	userController := cli.NewUserController(userUC)

	c.commands = []Command{
		{Name: "create-user", NameRu: "создать юзера", Handler: userController.Create},
		{Name: "delete-user", NameRu: "удалить юзера", Handler: userController.Delete},
		{Name: "all-users", NameRu: "список юзеров", Handler: userController.GetAll},
	}
}
