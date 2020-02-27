package main

import (
	"fmt"

	"github.com/the-maldridge/locksmith/internal/caller"
	"github.com/hashicorp/go-hclog"
)

var (
	wgLogger = hclog.New(&hclog.LoggerOptions{
		Name:  "WireGuard-Client-Test",
		Level: hclog.LevelFromString("DEBUG"),
	})
)

// This function receives a string and prints it with a border around it.
func PrintBorder(title string) {
	for i := 0; i < 79; i++ {
		fmt.Print("*")
	}
	fmt.Println("\n" + title)
	for i := 0; i < 79; i++ {
		fmt.Print("*")
	}
}

func main() {
	wgLogger.Info("Beginning Client type test.\n")

	// Make Client
	wgClient := caller.New()

	fmt.Println(wgClient, "\n")

	// Check helper functions

	wgClient.SetCompany("Big-Company")
	wgClient.UpdateConfiguration()
	fmt.Println(wgClient.GetCompany())

	fmt.Println(wgClient, "\n")
	PrintBorder("Testing helper functions")
	fmt.Println("\nPublic Key: ", wgClient.GetPublicKey())
	fmt.Println("Operating System: ", wgClient.GetOS())
	fmt.Println("Status: ", wgClient.GetStatus())
	fmt.Println("Changing status...")
	wgClient.SetStatus(caller.Active)
	fmt.Println("Status: ", wgClient.GetStatus())
	PrintBorder("Test Completed")
}
