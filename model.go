package main

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func NewModel(store *Store, charLimit int) model {
	notes, err := store.GetNotes()
	if err != nil {
		tea.Printf("Unable to get notes: %v", err)
		tea.Quit()
	}

	ta := initTextArea(charLimit)

	return model{
		store:     store,
		state:     listView,
		textarea:  ta,
		textinput: textinput.New(),
		notes:     notes,
	}
}

func initTextArea(limit int) textarea.Model {
	taTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a2009b"))
	taBaseStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#19d6c9")).Padding(0, 1)

	ta := textarea.New()
	ta.CharLimit = limit
	ta.ShowLineNumbers = false
	ta.Placeholder = "Say something, say something"

	taStyles := ta.Styles()
	taStyles.Focused.Base = taBaseStyle
	taStyles.Blurred.Base = taBaseStyle
	taStyles.Focused.Text = taTextStyle
	taStyles.Blurred.Text = taTextStyle

	ta.SetStyles(taStyles)

	return ta
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

		switch key {
		case "ctrl+esc":
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
						tea.Printf("Error saving note: %v", err)
					}

					m.notes, err = m.store.GetNotes()
					if err != nil {
						tea.Printf("Error getting notes: %v", err)
					}

					m.currNote = Note{}
					m.state = listView
				} else {
					// timer nonsense for fun
					m.textarea.SetValue("Type something you turd!")
					m.textarea.Focus()
					m.textarea.CursorEnd()
				}
			case "ctrl+c":
				cmds = append(cmds, tea.SetClipboard(m.textarea.SelectedText()))
			case "ctrl+x":
				// more stuff to do after copy so above method doesn't work
				if m.textarea.HasSelection() {
					selected := m.textarea.SelectedText()
					cmds = append(cmds, tea.SetClipboard(m.textarea.SelectedText()))
					ns := strings.Replace(m.textarea.Value(), selected, "", 1)
					m.textarea.SetValue(ns) // do something
					m.textarea.CursorEnd()
				}
			case "esc":
				m.state = listView
			}
		}
	}

	return m, tea.Batch(cmds...)
}
