package commander

import (
	"context"
	"space-go/internal/client"
	"space-go/internal/commander"
	"space-go/internal/commander/strategies/simple"
)

func main() {
	baseURL := "http://localhost:3333"
	token := "660c35366abee660c35366abf1"

	ctx := context.Background()
	apiClient := client.NewClient(baseURL, token)

	cmd := &commander.Commander{
		API:      apiClient,
		Strategy: simple.New(),
		State:    &commander.State{},
		Status:   commander.Running,
		Packer:   commander.DumboPacker{},
	}
	err := cmd.Run(ctx)
	if err != nil {
		panic(err)
	}
}
