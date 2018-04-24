package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	destDir, prevDir, unknownDir, mediainfo, renameString string
	movieExts, picExts                                    []string
	debug, noRenameDest, renameSource                     bool
	byteCount                                             int64
)

// object refers to the file we are examining in walkFunc
type object struct {
	FullSourcePath string
	SourceName     string
	SourcePath     string
	SourceSize     int64
	FullDestPath   string
	DestName       string
	DestPath       string
	FH             *os.File
}

func main() {

	var startDir, pics, movs string

	flag.StringVar(&startDir, "startDir", "", "Start traversing from here")
	flag.StringVar(&destDir, "destDir", "", "Put renamed files here")
	flag.StringVar(&unknownDir, "unknownDir", "", "Where to put files with unknown metadata")
	flag.StringVar(&pics, "picExts", "jpg,jpeg,gif,png,aae,tif,thm", "Comma delimited list of picture extensions")
	flag.StringVar(&movs, "movExts", "mov,mp4,avi,mod,m4a,m4v,lrv", "Comma delimited list of movie extensions")
	flag.StringVar(&mediainfo, "mediainfo", "/usr/bin/mediainfo", "Path to mediainfo binary")
	flag.BoolVar(&debug, "debug", false, "Set for more logging")
	flag.BoolVar(&noRenameDest, "noRenameDest", false, "Do not rename files on destination - keep existing source filename")
	flag.BoolVar(&renameSource, "renameSource", false, "Rename source file after successful copy")
	flag.Parse()

	if len(startDir) == 0 || len(destDir) == 0 {
		log.Fatalf("You did not specify either a starting or destination dir")
	}

	_, err := os.Stat(mediainfo)
	if err != nil {
		log.Fatalf("Did not find mediainfo at: %s", mediainfo)
	}

	// if we want to rename source file after copying, define our suffix
	if renameSource {
		renameString = fmt.Sprintf("copied-%d", time.Now().Unix())
	}

	// set this to our destination dir initially - it will be modified the next time we can successfully extract the time/date from a file
	prevDir = destDir

	picExts = parseExtensions(pics)
	movieExts = parseExtensions(movs)

	t1 := time.Now()

	err = filepath.Walk(startDir, walkFunc)
	if err != nil {
		log.Fatalf("filepath.Walk returned error: %v.\n", err)
	}

	throughput := float64(byteCount) / time.Since(t1).Seconds()
	log.Printf("Finished script. Copied %d bytes in %v. %.2f bytes/sec", byteCount, time.Since(t1), throughput)
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

	o := &object{
		SourceName:     filepath.Base(path),
		SourcePath:     filepath.Dir(path),
		SourceSize:     info.Size(),
		FullSourcePath: path,
	}

	fields := strings.Split(o.SourceName, `.`)
	extension := fields[len(fields)-1]
	lowerExtension := strings.ToLower(extension)

	for _, e := range picExts {
		if lowerExtension == e {
			examinePic(o)
			return nil
		}
	}

	for _, e := range movieExts {
		if lowerExtension == e {
			examineMov(o)
			return nil
		}
	}

	if debug {
		log.Printf("Not doing anything with %s - unknown extension", path)
	}

	return nil
}
