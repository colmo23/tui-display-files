"""TUI file browser built with Textual — Python port of the Go Bubble Tea app."""

import argparse
import logging
import os
import sys

from textual.app import App, ComposeResult
from textual.binding import Binding
from textual.containers import VerticalScroll
from textual.widgets import Header, Footer, ListView, ListItem, Static


from filesystem import read_dir, read_file_content


class FileItem(ListItem):
    """A list item representing a file or directory entry."""

    def __init__(self, name: str, desc: str, max_name_len: int) -> None:
        super().__init__()
        self.file_name = name
        self.desc = desc
        self.max_name_len = max_name_len

    def compose(self) -> ComposeResult:
        label = f"{self.file_name:<{self.max_name_len}}   {self.desc}"
        yield Static(label)


class FileBrowserApp(App):
    """A TUI file browser with directory navigation and file viewing."""

    TITLE = ""
    CSS = """
    Screen {
        background: $surface;
    }
    ListView {
        background: $surface;
    }
    ListView > ListItem {
        padding: 0 2;
    }
    ListView > ListItem.--highlight {
        background: $surface;
    }
    ListView > ListItem.--highlight > Static {
        color: $warning;
    }
    #file-content-title {
        color: $accent;
        text-style: bold;
        padding: 0 1;
    }
    #file-content-scroll {
        height: 1fr;
    }
    #file-content-body {
        padding: 0 2;
    }
    """

    BINDINGS = [
        Binding("q", "quit", "Quit"),
        Binding("escape", "go_back", "Back", show=True),
    ]

    def __init__(self, initial_dir: str) -> None:
        super().__init__()
        self.current_dir = os.path.abspath(initial_dir)
        self.view_state = "file_list"  # or "file_content"
        self.selected_file = ""

    def compose(self) -> ComposeResult:
        yield Header()
        yield ListView(id="file-list")
        yield Static("", id="file-content-title")
        yield VerticalScroll(Static("", id="file-content-body"), id="file-content-scroll")
        yield Footer()

    def on_mount(self) -> None:
        self._load_directory()
        self._show_file_list()

    def _load_directory(self) -> None:
        """Load directory entries into the ListView."""
        list_view = self.query_one("#file-list", ListView)
        list_view.clear()

        entries = read_dir(self.current_dir)

        # Build items list with ".." entry when not at root
        items = []
        abs_path = os.path.abspath(self.current_dir)
        parent = os.path.dirname(abs_path)
        if abs_path != parent:
            items.append(("..", "directory"))

        for entry in entries:
            if entry.is_dir():
                desc = "directory"
            else:
                try:
                    size = entry.stat().st_size
                    desc = f"{size} bytes"
                except OSError:
                    desc = ""
            items.append((entry.name, desc))

        # Calculate max name length for alignment
        max_name_len = max((len(name) for name, _ in items), default=0)

        for name, desc in items:
            list_view.append(FileItem(name, desc, max_name_len))

        self.title = self.current_dir

    def _show_file_list(self) -> None:
        """Switch to file list view."""
        self.view_state = "file_list"
        self.query_one("#file-list").display = True
        self.query_one("#file-content-title").display = False
        self.query_one("#file-content-scroll").display = False

    def _show_file_content(self, filename: str, content: str) -> None:
        """Switch to file content view."""
        self.view_state = "file_content"
        self.selected_file = filename
        self.query_one("#file-list").display = False
        title_widget = self.query_one("#file-content-title", Static)
        title_widget.update(filename)
        title_widget.display = True
        body_widget = self.query_one("#file-content-body", Static)
        body_widget.update(content)
        scroll = self.query_one("#file-content-scroll")
        scroll.display = True
        scroll.scroll_home(animate=False)

    def on_list_view_selected(self, event: ListView.Selected) -> None:
        """Handle Enter on a list item."""
        item = event.item
        if not isinstance(item, FileItem):
            return

        if item.desc == "directory":
            # Navigate into directory
            if item.file_name == "..":
                new_path = os.path.dirname(self.current_dir)
            else:
                new_path = os.path.join(self.current_dir, item.file_name)
            self.current_dir = os.path.abspath(new_path)
            logging.debug("Navigating to %s", self.current_dir)
            self._load_directory()
        else:
            # Open file content
            file_path = os.path.join(self.current_dir, item.file_name)
            try:
                content = read_file_content(file_path)
            except Exception as e:
                logging.error("Error reading file: %s", e)
                content = f"Error reading file: {e}"
            self._show_file_content(item.file_name, content)

    def action_go_back(self) -> None:
        """Return to file list from file content view."""
        if self.view_state == "file_content":
            self._show_file_list()


def main():
    parser = argparse.ArgumentParser(description="TUI file browser")
    parser.add_argument("--dir", default=".", help="Directory to display")
    parser.add_argument("--log", default="", dest="log_file", help="Path to log file")
    args = parser.parse_args()

    if args.log_file:
        logging.basicConfig(
            filename=args.log_file,
            level=logging.DEBUG,
            format="%(asctime)s %(levelname)s %(message)s",
        )
        logging.info("Logging enabled.")

    app = FileBrowserApp(args.dir)
    app.run()


if __name__ == "__main__":
    main()
