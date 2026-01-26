package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type viewState int

const (
	fileListView viewState = iota
	fileContentView
)

// item represents a file or directory in the list.
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type mainModel struct {
	state        viewState
	fileList     list.Model
	fileContent  viewport.Model
	selectedFile string
	ready        bool
	width, height int
	err          error
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.fileList.SetSize(msg.Width, msg.Height)
		m.fileContent.Width = msg.Width
		m.fileContent.Height = msg.Height - 1 // Reserve space for header
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case fileListView:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				selectedItem := m.fileList.SelectedItem().(item)
				if selectedItem.desc != "directory" {
					content, err := readFileContent(selectedItem.title)
					if err != nil {
						m.err = err
						return m, nil
					}
					m.fileContent.SetContent(content)
					m.selectedFile = selectedItem.title
					m.state = fileContentView
				}
			}
		case fileContentView:
			switch msg.String() {
			case "esc":
				m.state = fileListView
				return m, nil
			}
		}
	}

	if m.state == fileListView {
		m.fileList, cmd = m.fileList.Update(msg)
	} else {
		m.fileContent, cmd = m.fileContent.Update(msg)
	}

	return m, cmd
}

func (m mainModel) View() string {
	if !m.ready {
		return "Initializing..."
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	if m.state == fileContentView {
		header := headerStyle.Render(m.selectedFile)
		return fmt.Sprintf("%s\n%s", header, m.fileContent.View())
	}
	return docStyle.Render(m.fileList.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func initialModel() mainModel {
	items := []list.Item{}
	files, err := readDir(".")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		var desc string
		if file.IsDir() {
			desc = "directory"
		} else {
			info, err := file.Info()
			if err == nil {
				desc = fmt.Sprintf("%d bytes", info.Size())
			}
		}
		items = append(items, item{title: file.Name(), desc: desc})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a file"
	l.SetShowHelp(true)

	vp := viewport.New(80, 24)

	return mainModel{
		state:       fileListView,
		fileList:    l,
		fileContent: vp,
	}
}
