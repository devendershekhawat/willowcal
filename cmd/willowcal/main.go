package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devendershekhawat/teambiscuit/internal/config"
)

func main() {
    // Parse CLI args
    if len(os.Args) < 2 {
        log.Fatal("Usage: teambiscuit init <config.yaml>")
    }
    
    configPath := os.Args[1]
    
    // Parse config
    cfg, err := config.ParseConfigFile(configPath)
    if err != nil {
        log.Fatalf("❌ Failed to parse config: %v", err)
    }
    
	fmt.Println(cfg)
}