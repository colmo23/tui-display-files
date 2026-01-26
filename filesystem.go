package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

// readDir reads the contents of a directory and returns a slice of fs.DirEntry,
// filtering out special directories like .git, .gemini, and hidden files.
func readDir(dirPath string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var files []fs.DirEntry
	for _, entry := range entries {
		// Ignore special directories and hidden files
		if entry.IsDir() && (entry.Name() == ".git" || entry.Name() == ".gemini") {
			continue
		}
		if !entry.IsDir() && filepath.Base(entry.Name()) == ".DS_Store" {
			continue
		}
		// Ignore hidden files and directories (those starting with a dot)
		if len(entry.Name()) > 0 && entry.Name()[0] == '.' && entry.Name() != "." && entry.Name() != ".." {
			continue
		}
		files = append(files, entry)
	}
	return files, nil
}

// readFileContent reads the entire content of a file and returns it as a string.
// It handles potential errors during file reading.
func readFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
