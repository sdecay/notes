package main

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	state     uint
	store     *Store
	notes     []Note
	currNote  Note
	listIndex int
	textarea  textarea.Model
	textinput textinput.Model
}

func NewModel(store *Store) model {
	notes, err := store.GetNotes()
	if err != nil {
		tea.Printf("unable to get notes: %v", err)
		tea.Quit()
	}

	return model{
		store:     store,
		state:     listView,
		textarea:  textarea.New(),
		textinput: textinput.New(),
		notes:     notes,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()

		// global ctrl+c to break out
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch m.state {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "n":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.currNote = Note{}
				m.state = titleView
				// show input
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}
			case "enter", "l":
				m.currNote = m.notes[m.listIndex]

				m.state = bodyView
				m.textarea.SetValue(m.currNote.Body)
				m.textarea.Focus()
				m.textarea.CursorEnd()
			case "ctrl+d":
				err := m.store.DeleteNote(&m.notes, m.listIndex)
				if err != nil {
					tea.Printf("error deleting note: %v\n", err)
					tea.Quit()
				}

				// m.notes = notes
				m.listIndex = 0
				m.state = listView
			}

		case titleView:
			switch key {
			case "enter":
				title := m.textinput.Value()

				if title != "" {
					m.currNote.Title = title

					m.state = bodyView
					m.textarea.SetValue("")
					m.textarea.Focus()
					m.textarea.CursorEnd()
				}
			case "esc":
				m.state = listView
			}

		case bodyView:
			switch key {
			case "ctrl+s":
				m.currNote.Body = m.textarea.Value()

				if m.currNote.Body != "" {
					var err error
					if err = m.store.SaveNote(m.currNote); err != nil {
						// proper error handling
						return m, tea.Quit
					}

					m.notes, err = m.store.GetNotes()
					if err != nil {
						// proper error handling
						return m, tea.Quit
					}

					m.currNote = Note{}
					m.state = listView
				}
			case "esc":
				m.state = listView
			}
		}
	}

	return m, tea.Batch(cmds...)
}
