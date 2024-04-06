package main

import (
	"context"
	"space-go/internal/client"
	"space-go/internal/commander"
	"space-go/internal/commander/packer"
	"space-go/internal/commander/strategies/simple"
	"time"
)

func main() {
	//baseURL := "http://localhost:3333"
	baseURL := "https://datsedenspace.datsteam.dev"
	token := "660c35366abee660c35366abf1"

	time.Sleep(3 * time.Second)

	ctx := context.Background()
	apiClient := client.NewClient(baseURL, token)

	cmd := &commander.Commander{
		API:      apiClient,
		Strategy: simple.New(),
		State:    &commander.State{},
		Status:   commander.Running,
		Packer:   packer.DuploPacker{},
	}

	err := cmd.Run(ctx)
	if err != nil {
		panic(err)
	}
}
