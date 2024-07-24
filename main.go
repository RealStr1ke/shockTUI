package main;

import (
	"shockTUI/src"
);

func main() {
	files, err := src.GetPages("pages");
	if err != nil {
		panic(err);
	}
	for _, file := range files {
		// Omit the file extensions and print the file names
		println(file.Name()[:len(file.Name())-3]);
	}
}