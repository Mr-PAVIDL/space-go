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
	Packer   Packer
}

func (commander *Commander) Run(ctx context.Context) error {
	if err := commander.init(ctx); err != nil {
		return fmt.Errorf("failed to init: %w", err)
	}

	commander.Status = Running
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

func (commander *Commander) init(ctx context.Context) error {
	universeResponse, err := commander.API.GetUniverse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get universe: %w", err)
	}

	state := &State{
		Planet:    universeResponse.Ship.Planet,
		FuelUsed:  universeResponse.Ship.FuelUsed,
		CapacityX: universeResponse.Ship.CapacityX,
		CapacityY: universeResponse.Ship.CapacityY,
		Garbage:   universeResponse.Ship.Garbage,
	}

	for _, garbage := range universeResponse.Ship.Garbage {
		state.GarbagePiecesCollected += 1
		state.GarbageCellsCollected += len(garbage)
	}

	state.Universe = NewUniverse(universeResponse.Universe)
	state.Universe.Planets[state.Planet.Name] = state.Planet

	commander.State = state

	return nil
}
