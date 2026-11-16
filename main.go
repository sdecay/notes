package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {
	store := new(Store)
	if err := store.Init(); err != nil {
		log.Fatalf("Unable to init: %v", err)
	}

	m := NewModel(store)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("some shit went wonky: %v", err)
	}
}
