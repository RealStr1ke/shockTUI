package main;

import (
	"shockTUI/src"
);

func printPages(dir string) error {
	pages, err := src.GetPages(dir);
	if err != nil {
		return err;
	}

	for page, content := range pages {
		println(page);
		println("=======");
		println(content);
		println("\n");1
	}

	return nil;
}

func main() {
	// files, err := GetPages("pages");
	// if err != nil {
	// 	panic(err);
	// }
	// for _, file := range files {
	// 	// Omit the file extensions and print the file names
	// 	println(file.Name()[:len(file.Name())-3]);
	// }

	printPages("pages");
}