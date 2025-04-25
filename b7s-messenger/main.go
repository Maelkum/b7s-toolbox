package main

import (
	"errors"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"

	"github.com/Maelkum/b7s-toolbox/b7s-messenger/messenger"
)

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	var (
		flagNode string
	)

	pflag.StringVarP(&flagNode, "node", "n", "", "multiaddress of the node to connect to")

	pflag.Parse()

	if flagNode == "" {
		return errors.New("node address is required")
	}

	model := model{
		choices:  messenger.SupportedMessages(),
		selected: make(map[int]struct{}),
	}

	program := tea.NewProgram(model)
	_, err := program.Run()
	if err != nil {
		return fmt.Errorf("error running messenger: %w", err)
	}

	return nil
}
