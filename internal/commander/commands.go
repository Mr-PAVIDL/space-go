package commander

import (
	"context"
	"fmt"
	"space-go/internal/model"
	"strings"
	"time"
)

type Command interface {
	Execute(ctx context.Context, commander *Commander) error
}

type TravelCommand struct {
	Path []string
}

func Travel(planets ...string) TravelCommand {
	return TravelCommand{planets}
}

func (cmd TravelCommand) Execute(ctx context.Context, commander *Commander) error {
	response, err := commander.API.Travel(ctx, model.TravelRequest{
		Planets: cmd.Path,
	})

	if err != nil {
		return fmt.Errorf("failed to make a travel api call: %w", err)
	}

	// TODO: apply diffs tp the state

	return nil
}

func (cmd TravelCommand) String() string {
	return formatCommand("travel", strings.Join(cmd.Path, " -> "))
}

type CollectCommand struct {
}

func Collect() CollectCommand {
	return CollectCommand{}
}

func (cmd CollectCommand) Execute(ctx context.Context, commander *Commander) error {
	// TODO: use collector from commander to optimally collect garbage

	return nil
}

func (cmd CollectCommand) String() string {
	return formatCommand("collect", "collecting garbage...")
}

type GotoCommand struct {
	Destination string
}

func GoTo(planet string) GotoCommand {
	return GotoCommand{}
}

func (cmd GotoCommand) Execute(ctx context.Context, commander *Commander) error {
	// TODO: use collector from commander to optimally collect garbage

	return nil
}

func (cmd GotoCommand) String() string {
	return formatCommand("goto", cmd.Destination)
}

type IdleCommand struct {
	Duration time.Duration
}

func Idle(duration time.Duration) IdleCommand {
	return IdleCommand{}
}

func (cmd IdleCommand) Execute(ctx context.Context, commander *Commander) error {
	time.Sleep(cmd.Duration)

	return nil
}

func (cmd IdleCommand) String() string {
	return formatCommand("idle", cmd.Duration.String())
}

type SequentialCommand []Command

func Sequential(commands ...Command) SequentialCommand {
	return commands
}

func (cmd SequentialCommand) Execute(ctx context.Context, commander *Commander) error {
	for _, subCmd := range cmd {
		if err := commander.Execute(ctx, subCmd); err != nil {
			return err
		}
	}

	return nil
}

func (cmd SequentialCommand) String() string {
	return formatCommand("sequential", fmt.Sprintf("%d", len(cmd)))
}

func formatCommand(name string, message string) string {
	return fmt.Sprintf("[%s]: %s", name, message)
}
