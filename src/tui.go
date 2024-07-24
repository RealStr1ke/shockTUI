package src;

import (
	"os"
	"path/filepath"
);

func GetPages(dir string) (map[string]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err;
	}
	
	pages := make(map[string]string)
	for _, file := range files {
		// pages = append(pages, file.Name()[:len(file.Name())-3]);
		content, err := os.ReadFile(filepath.Join(dir, file.Name()));
		if err != nil {
			return nil, err;
		}
		pages[file.Name()[:len(file.Name())-3]] = string(content);
	}

	return pages, nil;
}
