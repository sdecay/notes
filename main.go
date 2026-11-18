package main

import (
	tea "charm.land/bubbletea/v2"
)

func main() {
	taLimit := 1024

	store := new(Store)
	if err := store.Init(); err != nil {
		tea.Printf("Unable to init: %v", err)
		tea.Quit()
	}

	m := NewModel(store, taLimit)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		tea.Printf("Some shit went wonky: %v", err)
		tea.Quit()
	}
}
