package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go-tui-library/modules/core"
	"go-tui-library/modules/renderer"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <config.json>")
		os.Exit(1)
	}

	configFile := os.Args[1]
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	var config core.UIConfig
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	model := renderer.NewModel(config)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
