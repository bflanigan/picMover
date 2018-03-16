package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	destDir, prevDir, movieDir, unknownDir string
	movieExts, picExts                     []string
)

func main() {

	var startDir, pics, movs string

	flag.StringVar(&startDir, "startDir", "", "Start traversing from here")
	flag.StringVar(&destDir, "destDir", "", "Put renamed files here")
	flag.StringVar(&movieDir, "movieDir", "", "Where to dump movies")
	flag.StringVar(&unknownDir, "unknownDir", "", "Where to put files with unknown metadata")
	flag.StringVar(&pics, "picExts", "jpg,gif,png,aae,tif", "Comma delimited list of picture extensions")
	flag.StringVar(&movs, "movExts", "mov,mp4,avi,mod,m4a,m4v", "Comma delimited list of picture extensions")
	flag.Parse()

	if len(startDir) == 0 || len(destDir) == 0 {
		log.Fatalf("You did not specify either a starting or destination dir")
	}

	prevDir = destDir

	picExts = parseExtensions(pics)
	movieExts = parseExtensions(movs)

	err := filepath.Walk(startDir, walkFunc)
	if err != nil {
		log.Fatalf("filepath.Walk returned error: %v.\n", err)
	}

}

func walkFunc(path string, info os.FileInfo, err error) error {

	if err != nil {
		log.Printf("Got error entering walkFunc for file %s - skipping, error: %v", path, err)
		return nil
	}

	if info.IsDir() {
		log.Printf("Skipping dir: %s", path)
		return nil
	}

	fields := strings.Split(info.Name(), `.`)
	extension := fields[len(fields)-1]
	lowerExtension := strings.ToLower(extension)

	for _, e := range picExts {
		if lowerExtension == e {
			examinePic(path)
			return nil
		}
	}

	for _, e := range movieExts {
		if lowerExtension == e {
			examineMov(path)
			return nil
		}
	}

	log.Printf("Not doing anything with %s - unknown extension", path)

	return nil
}
