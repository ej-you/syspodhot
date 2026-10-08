// Package climanager provides command line manager service.
package climanager

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"syspodhot/config"
	"syspodhot/internal/pkg/db"
)

// Handler represents a command handler.
type Handler func() error

// Command represents a command with name and handler.
type Command struct {
	Name    string
	NameRu  string
	Handler Handler
}

// CLIManager represents a command line manager.
type CLIManager struct {
	cfg      *config.Config
	commands []Command

	started chan struct{}
	exited  chan struct{}
}

// New returns a new instance of [CLIManager].
func New(cfg *config.Config, dbStorage *db.DB) (*CLIManager, error) {
	manager := &CLIManager{
		started: make(chan struct{}),
		exited:  make(chan struct{}),
		cfg:     cfg,
	}
	// register all commands
	manager.registerCommands(dbStorage)
	return manager, nil
}

// Started returns chan that signals the service's readiness.
func (c *CLIManager) Started() <-chan struct{} {
	return c.started
}

// Exited returns chan that signals that service's exited.
func (c *CLIManager) Exited() <-chan struct{} {
	return c.exited
}

// StartWithShutdown starts cmd manager and waits for
// context is done for gracefully shutdown it.
// This method is blocking.
func (c *CLIManager) Run(ctx context.Context, errChan chan<- error) {
	slog.Info("start cli manager service...")

	// start cli manager
	go func() {
		for {
			fmt.Println(c.msgAvailableCommands())
			if err := c.serve(); err != nil {
				errChan <- fmt.Errorf("cli manager: %w", err)
				return
			}
		}
	}()

	close(c.started)
	defer close(c.exited)
	// wait for context
	<-ctx.Done()
}

// msgAvailableCommands returns a string with list of all cli manager commands.
func (c *CLIManager) msgAvailableCommands() string {
	var builder strings.Builder

	builder.WriteString("\nДоступные команды:")
	// collect available commands slice to string
	for idx, command := range c.commands {
		builder.WriteString("\n\t ")
		builder.WriteString(strconv.Itoa(idx + 1))
		builder.WriteString(". ")
		builder.WriteString(command.NameRu)
	}
	return builder.String()
}

// serve represents a one iteration of infinity loop to serve registered command.
func (c *CLIManager) serve() error {
	var (
		cmdIdx int
		cmd    *Command
		err    error
	)
	for cmd == nil {
		if cmdIdx = c.scanCmdIdx(); cmdIdx == -1 {
			continue
		}
		cmd = &c.commands[cmdIdx]
	}
	// exec command
	if err = cmd.Handler(); err != nil {
		return fmt.Errorf("command %s: %w", cmd.Name, err)
	}
	return nil
}

// scanCmdIdx asks for index of command and scans it.
// Returns -1 if invalid command index was given.
func (c *CLIManager) scanCmdIdx() int {
	var cmdIdxRaw string

	// ask for command
	fmt.Print("Введите номер команды (или используйте Ctrl+C для выхода): ")
	fmt.Scan(&cmdIdxRaw)

	cmdIdx, err := strconv.Atoi(cmdIdxRaw)
	if err != nil {
		fmt.Println("Введите корректный номер!")
		return -1
	}
	cmdIdx-- // starts from zero
	if cmdIdx < 0 || cmdIdx >= len(c.commands) {
		fmt.Println("Команда по номеру не найдена!")
		return -1
	}
	return cmdIdx
}
