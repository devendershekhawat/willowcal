package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devendershekhawat/teambiscuit/internal/commands"
)

func main() {
	// Parse CLI args
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Execute command
	var err error
	switch command {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("❌ Missing config file path")
			printUsage()
			os.Exit(1)
		}
		configPath := os.Args[2]
		err = commands.InitCommand(configPath)
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("❌ Missing config file path")
			printUsage()
			os.Exit(1)
		}
		configPath := os.Args[2]
		err = commands.RunCommand(configPath)
	case "server":
		port := "8080"
		workspaceDir := "./workspace"
		staticDir := "./web/dist"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}
		if len(os.Args) > 3 {
			workspaceDir = os.Args[3]
		}
		if len(os.Args) > 4 {
			staticDir = os.Args[4]
		}
		err = commands.ServerCommand(port, workspaceDir, staticDir)
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
	fmt.Println("  willowcal <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init <config.yaml>           Clone repositories and run setup commands")
	fmt.Println("  run <config.yaml>            Start services (clone missing repos if needed)")
	fmt.Println("  server [port] [workspace]    Start WebSocket server (default port: 8080)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  willowcal init config.yaml")
	fmt.Println("  willowcal run config.yaml")
	fmt.Println("  willowcal server")
	fmt.Println("  willowcal server 3000 ./my-workspace")
}