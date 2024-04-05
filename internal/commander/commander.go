package commander

import (
	"context"
	"fmt"
	"log/slog"
	"space-go/internal/client"
)

type Commander struct {
	API      *client.DatsEdenSpaceClient
	Strategy Strategy
	State    *State
	Status   Status
}

func (commander *Commander) Run(ctx context.Context) error {
	for commander.Status == Running {
		cmd := commander.Strategy.Next(ctx, commander.State)
		if err := commander.Execute(ctx, cmd); err != nil {
			slog.Error("failed to execute a command", slog.String("error", err.Error()))
			continue
		}

	}

	return nil
}

func (commander *Commander) Execute(ctx context.Context, cmd Command) error {
	if stringer, ok := cmd.(fmt.Stringer); ok {
		slog.Info("running a command", slog.String("cmd", stringer.String()))
	}
	return cmd.Execute(ctx, commander)
}
