package main

import (
	"context"
	"fmt"
	"space-go/internal/client"
)

func main() {
	baseURL := "https://datsedenspace.datsteam.dev"
	token := "660c35366abee660c35366abf1"

	ctx := context.Background()
	apiClient := client.NewClient(baseURL, token)

	player, err := apiClient.GetUniverse(ctx)
	if err != nil {
		fmt.Printf("Error getting universe: %v\n", err)
		return
	}
	fmt.Printf("Player: %v\n", player)

}
