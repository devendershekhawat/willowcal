package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devendershekhawat/teambiscuit/internal/commands"
)

func main() {
	// Parse CLI args
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	configPath := os.Args[2]

	// Execute command
	var err error
	switch command {
	case "init":
		err = commands.InitCommand(configPath)
	case "run":
		err = commands.RunCommand(configPath)
	default:
		fmt.Printf("❌ Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("❌ %v", err)
	}
}

func printUsage() {
	fmt.Println("willowcal - Repository orchestration tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  willowcal <command> <config.yaml>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init    Clone repositories and run setup commands")
	fmt.Println("  run     Start services (clone missing repos if needed)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  willowcal init config.yaml")
	fmt.Println("  willowcal run config.yaml")
}