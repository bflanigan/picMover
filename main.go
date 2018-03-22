package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	//destDir, prevDir, movieDir, unknownDir, mediainfo string
	destDir, prevDir, unknownDir, mediainfo string
	movieExts, picExts                      []string
	debug, noRename                         bool
	byteCount                               int64
)

func main() {

	var startDir, pics, movs string

	flag.StringVar(&startDir, "startDir", "", "Start traversing from here")
	flag.StringVar(&destDir, "destDir", "", "Put renamed files here")
	//flag.StringVar(&movieDir, "movieDir", "", "Where to dump movies")
	flag.StringVar(&unknownDir, "unknownDir", "", "Where to put files with unknown metadata")
	flag.StringVar(&pics, "picExts", "jpg,gif,png,aae,tif,thm", "Comma delimited list of picture extensions")
	flag.StringVar(&movs, "movExts", "mov,mp4,avi,mod,m4a,m4v,lrv", "Comma delimited list of picture extensions")
	flag.StringVar(&mediainfo, "mediainfo", "/usr/bin/mediainfo", "Path to mediainfo binary")
	flag.BoolVar(&debug, "debug", false, "Set for more logging")
	flag.BoolVar(&noRename, "noRename", false, "Do not rename files - keep existing filename")
	flag.Parse()

	if len(startDir) == 0 || len(destDir) == 0 {
		log.Fatalf("You did not specify either a starting or destination dir")
	}

	_, err := os.Stat(mediainfo)
	if err != nil {
		log.Fatalf("Did not find mediainfo at: %s", mediainfo)
	}

	prevDir = destDir

	picExts = parseExtensions(pics)
	movieExts = parseExtensions(movs)

	t1 := time.Now()

	err = filepath.Walk(startDir, walkFunc)
	if err != nil {
		log.Fatalf("filepath.Walk returned error: %v.\n", err)
	}

	secs := time.Since(t1).Seconds()
	throughput := float64(byteCount) / secs

	log.Printf("Finished script. Copied %d in %v. %f bytes/sec", byteCount, time.Since(t1), throughput)

}

func walkFunc(path string, info os.FileInfo, err error) error {

	if err != nil {
		log.Printf("Got error entering walkFunc for file %s - skipping, error: %v", path, err)
		return nil
	}

	if info.IsDir() {
		if debug {
			log.Printf("Skipping dir: %s", path)
		}
		return nil
	}

	// ignore ._IMG files
	if info.Size() == 4096 && strings.HasPrefix(filepath.Base(path), "._") {
		if debug {
			log.Printf("Skipping file: %s", path)
		}
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

	if debug {
		log.Printf("Not doing anything with %s - unknown extension", path)
	}

	return nil
}
