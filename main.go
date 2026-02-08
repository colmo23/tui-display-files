package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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

// fileItemDelegate is a custom delegate for list items
type fileItemDelegate struct{}

func (d fileItemDelegate) Height() int {
	return 1
}

func (d fileItemDelegate) Spacing() int {
	return 0
}

func (d fileItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d fileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s   %s", i.Title(), i.Description())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(s...)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type mainModel struct {
	state         viewState
	fileList      list.Model
	fileContent   viewport.Model
	currentDir    string
	selectedFile  string
	ready         bool
	width, height int
	err           error
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
				// Handle directory traversal
				if selectedItem.desc == "directory" {
					var newPath string
					if selectedItem.title == ".." { // Go up one directory
						newPath = filepath.Dir(m.currentDir)
					} else { // Go into selected directory
						newPath = filepath.Join(m.currentDir, selectedItem.title)
					}

					// Update file list for new directory
					newItems, err := loadDirectoryItems(newPath)
					if err != nil {
						log.Printf("error loading directory items: %v", err)
						m.err = err
						return m, nil
					}
					m.fileList.SetItems(newItems)
					//			m.selectedFile = newItems[0].FilterValue()
					m.currentDir = newPath
					m.fileList.Title = fmt.Sprintf("Files in %s", m.currentDir)
					return m, nil
				} else { // It's a file, view its content
					content, err := readFileContent(filepath.Join(m.currentDir, selectedItem.title))
					if err != nil {
						log.Printf("error reading file content: %v", err)
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
	dirPtr := flag.String("dir", ".", "the directory to display")
	logFilePtr := flag.String("log", "", "path to log file")
	flag.Parse()

	if *logFilePtr != "" {
		f, err := os.OpenFile(*logFilePtr, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(f)
		log.Println("Logging enabled.")
	}

	p := tea.NewProgram(initialModel(*dirPtr), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func initialModel(initialDir string) mainModel {
	items, err := loadDirectoryItems(initialDir)
	if err != nil {
		log.Printf("error loading initial directory items: %v", err)
		panic(err)
	}

	// Create a new custom delegate
	delegate := fileItemDelegate{}
	l := list.New(items, delegate, 0, 0)
	l.Title = fmt.Sprintf("Files in %s", initialDir)
	l.SetShowHelp(true)

	vp := viewport.New(80, 24)

	return mainModel{
		state:       fileListView,
		fileList:    l,
		fileContent: vp,
		currentDir:  initialDir,
	}
}

func loadDirectoryItems(dirPath string) ([]list.Item, error) {
	entries, err := readDir(dirPath)
	if err != nil {
		return nil, err
	}

	items := []list.Item{}
	// Add ".." to go up a directory, if not at root
	if dirPath != "/" && dirPath != "." { // Check for root and current directory
		absPath, err := filepath.Abs(dirPath)
		if err == nil {
			parentDir := filepath.Dir(absPath)
			if absPath != parentDir { // Ensure we're not adding ".." at filesystem root
				items = append(items, item{title: "..", desc: "directory"})
			}
		}
	}

	for _, entry := range entries {
		var desc string
		if entry.IsDir() {
			desc = "directory"
		} else {
			info, err := entry.Info()
			if err == nil {
				desc = fmt.Sprintf("%d bytes", info.Size())
			}
		}
		items = append(items, item{title: entry.Name(), desc: desc})
	}
	return items, nil
}
