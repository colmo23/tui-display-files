# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
./build.sh            # builds binary: ./tui-display-files
go build -o tui-display-files  # equivalent
./tui-display-files              # run in current directory
./tui-display-files -dir /path   # run in specific directory
./tui-display-files -log debug.log  # enable logging to file
```

## Test

```bash
go test ./...         # run all tests
go test -run TestReadDir  # run a single test
```

## Architecture

This is a Go TUI file browser built with the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework (Elm architecture: Model → Update → View).

**Single model, two view states:** The app uses one `mainModel` struct with a `viewState` enum (`fileListView` / `fileContentView`) rather than separate models per view. State transitions happen via keybindings in `Update()`.

- `main.go` — Bubble Tea program: `mainModel` (state), `Update` (input handling + state transitions), `View` (rendering), `item` + `fileItemDelegate` (custom list item rendering), `loadDirectoryItems` (builds list items from directory entries). Entry point parses `-dir` and `-log` flags.
- `filesystem.go` — Filesystem abstraction: `readDir` (reads directory, filters hidden files/`.git`/`.gemini`/`.DS_Store`), `readFileContent` (reads file to string).
- `styles.go` — Lipgloss style definitions for the list and content views.
- `filesystem_test.go` — Tests for `readDir` and `readFileContent`.

**Key components from Charm libraries:**
- `bubbles/list.Model` — file list with filtering/navigation
- `bubbles/viewport.Model` — scrollable file content viewer
- `lipgloss` — terminal styling

**Navigation flow:** File list → Enter on directory navigates into it (with `..` support) → Enter on file opens content in viewport → Esc returns to file list.
