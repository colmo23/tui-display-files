import os


def read_dir(dir_path):
    """Read directory entries, filtering hidden files, .git, .gemini, .DS_Store."""
    entries = sorted(os.scandir(dir_path), key=lambda e: e.name)
    result = []
    for entry in entries:
        name = entry.name
        # Ignore .git and .gemini directories
        if entry.is_dir() and name in (".git", ".gemini"):
            continue
        # Ignore .DS_Store files
        if not entry.is_dir() and name == ".DS_Store":
            continue
        # Ignore hidden files/directories (starting with dot)
        if name.startswith(".") and name not in (".", ".."):
            continue
        result.append(entry)
    return result


def read_file_content(file_path):
    """Read the entire content of a file and return it as a string."""
    with open(file_path, "r") as f:
        return f.read()
