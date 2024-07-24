package src;

import (
	"os"
	"path/filepath"
);

func GetPages(dir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err;
	}
	
	// Return all the files that end in .md
	var pages []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			continue;
		}
		if filepath.Ext(file.Name()) == ".md" {
			pages = append(pages, file);
		}
	}

	return pages, nil;
}
