# Implementation Plan: TUI File Viewer

This document provides a step-by-step plan for building the TUI File Viewer application, based on the technical specification.

---

### Phase 1: Environment & Project Init

-   [x] **1.1:** Set up and verify the local Go development environment (`>=1.18`).
-   [x] **1.2:** Create a new project directory.
-   [x] **1.3:** Initialize the Go module with `go mod init <module_name>`.
-   [x] **1.4:** Add core dependencies using `go get`:
    -   [x] `github.com/charmbracelet/bubbletea`
    -   [x] `github.com/charmbracelet/bubbles`
    -   [x] `github.com/charmbracelet/lipgloss`
-   [x] **1.5:** Create the initial project structure (`main.go`).

### Phase 2: Backend/Data Layer (File System Interaction)

-   [x] **2.1:** Implement a function to read a list of files and directories from a given path.
    -   [x] Handle errors, such as permissions issues or non-existent directories.
-   [x] **2.2:** Implement a function to read the complete content of a specified file into a string.
    -   [x] Handle errors, such as the file being a binary or being unreadable.

### Phase 3: Core Logic & Features

-   [x] **3.1:** Define the `mainModel` and `fileViewModel` structs as specified in `TECH-SPEC.md`.
-   [x] **3.2:** Implement the state management for the `mainModel` (file list view).
    -   [x] **Init:** Load the initial list of files from the current directory.
    -   [x] **Update:** Handle keyboard input (Up/Down arrows) to move the selection in the list.
    -   [x] **Update:** Handle terminal resize events.
-   [x] **3.3:** Implement the state transition from `mainModel` to `fileViewModel`.
    -   [x] **Update:** On 'Enter' key press, capture the selected file path and instantiate the `fileViewModel`.
-   [x] **3.4:** Implement the state management for the `fileViewModel` (file content view).
    -   [x] **Init:** Load the content of the selected file.
    -   [x] **Update:** Handle keyboard input (PageUp/PageDown) to scroll the content.
    -   [x] **Update:** Handle terminal resize events.
    -   [x] **Update:** Handle 'Escape' key press to transition back to the `mainModel`.

### Phase 4: Frontend/UI (View Rendering)

-   [x] **4.1:** Implement the `View()` method for the `mainModel`.
    -   [x] Use the `bubbles/list` component to render the directory contents.
    -   [x] Use `lipgloss` to style the list, highlighting the selected item.
    -   [x] Add help text indicating basic controls.
-   [x] **4.2:** Implement the `View()` method for the `fileViewModel`.
    -   [x] Use the `bubbles/viewport` component to render the file's content.
    -   [x] Display the file path and scroll progress.
    -   [x] Add help text for scrolling and returning to the list.

### Phase 5: Testing & Deployment

-   [x] **5.1:** Write unit tests for the file system functions (Phase 2).
-   [x] **5.2:** Perform comprehensive manual testing of the entire user flow.
    -   [x] Test navigation, selection, and scrolling.
    -   [x] Test with different directory sizes and file types.
    -   [x] Test edge cases (empty directory, empty file, permission errors).
-   [x] **5.3:** Create a build script (`build.sh` or `Makefile`) for producing application binaries.
-   [x] **5.4:** Write a `README.md` file explaining what the application is, how to build it, and how to run it.
-   [x] **5.5:** Create a git tag for the `v1.0.0` release.
