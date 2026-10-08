// App binary starts full application.
package main

import (
	"log/slog"
	"os"

	"syspodhot/internal/app"
)

const (
	_argServer  = "serve"   // arg to run HTTP-server
	_argManager = "manager" // arg to run CLI-manager
)

func main() {
	var (
		err         error
		application *app.App
	)

	// parse program args (without program name)
	args := os.Args[1:]

	switch {
	case len(args) == 0 || args[0] == _argServer:
		application, err = app.NewServe()
	case args[0] == _argManager:
		application, err = app.NewCLI()
	}
	if err != nil {
		fatal(err)
	}
	if err := application.Run(); err != nil {
		fatal(err)
	}
}

// fatal logs error and exit with code 1.
func fatal(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}
