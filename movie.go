package main

import (
	"log"
	"os"
	"path/filepath"
)

func examineMov(path string) {

	var destinationDir string
	if len(movieDir) == 0 {
		destinationDir = prevDir
	} else {
		destinationDir = movieDir
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open sourcefile: %s due to error: %v", path, err)
	}
	defer f.Close()

	copyFile(f, filepath.Base(path), destinationDir)
}
