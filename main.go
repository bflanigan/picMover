package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var (
	startDir, destDir, prevDir string
)

func main() {

	var analyze bool
	var checkFile string

	flag.StringVar(&startDir, "startDir", "", "Start traversing from here")
	flag.StringVar(&destDir, "destDir", "", "Put renamed files here")
	flag.BoolVar(&analyze, "analyze", false, "Check for metadata on file")
	flag.StringVar(&checkFile, "checkfile", "", "File to examine")
	flag.Parse()

	prevDir = "/media/camera"

	if analyze {
		if len(checkFile) == 0 {
			log.Fatalf("You did not specify a file to examine")
		}

		f, err := os.Open(checkFile)
		if err != nil {
			log.Fatalf("Failed to open %s for reading due to error: %v", checkFile, err)
		}
		defer f.Close()

		x, err := exif.Decode(f)
		if err != nil {
			log.Fatalf("Failed to decode exif tag due to error: %v", err)
		}

		dateTaken, err := x.DateTime()
		if err != nil {
			log.Fatalf("Parsed exif tag but was unable to extract timestamp due to error: %v", err)
		}

		var strday, strhour, strmin, strsec string

		year, month, day := dateTaken.Date()
		hour, min, sec := dateTaken.Clock()

		stryear := strconv.Itoa(year)
		strmonth := numMonth(month)

		if day < 10 {
			strday = strconv.Itoa(day)
			strday = "0" + strday
		} else {
			strday = strconv.Itoa(day)
		}

		if hour < 10 {
			strhour = strconv.Itoa(hour)
			strhour = "0" + strhour
		} else {
			strhour = strconv.Itoa(hour)
		}

		if min < 10 {
			strmin = strconv.Itoa(min)
			strmin = "0" + strmin
		} else {
			strmin = strconv.Itoa(min)
		}

		if sec < 10 {
			strsec = strconv.Itoa(sec)
			strsec = "0" + strsec
		} else {
			strsec = strconv.Itoa(sec)
		}

		//origFilename := filepath.Base(checkFile)
		//newFilename := fmt.Sprintf("%s%s%s_%s%s%s_%s", stryear, strmonth, strday, strhour, strmin, strsec, origFilename)

		fmt.Printf("%s - %s%s%s_%s%s%s\n", checkFile, stryear, strmonth, strday, strhour, strmin, strsec)
		os.Exit(1)

	}

	if len(startDir) == 0 || len(destDir) == 0 {
		log.Fatalf("You did not specify either a starting or destination dir")
	}

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

	// copy AAE files
	if info.Size() < 8192 {
		copyInvalidExtension(path)
		return nil
	}

	fields := strings.Split(info.Name(), `.`)
	extension := fields[len(fields)-1]
	lowerExtension := strings.ToLower(extension)

	// get PNG, MOV, MP4
	if lowerExtension != "jpg" {
		copyInvalidExtension(path)
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %s for reading due to error: %v", path, err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		copyMissingEXIF(f)
		return nil
	}

	dateTaken, err := x.DateTime()
	if err != nil {
		copyMissingEXIF(f)
		return nil
	}

	var strday, strhour, strmin, strsec string

	year, month, day := dateTaken.Date()
	hour, min, sec := dateTaken.Clock()

	stryear := strconv.Itoa(year)
	strmonth := numMonth(month)

	if day < 10 {
		strday = strconv.Itoa(day)
		strday = "0" + strday
	} else {
		strday = strconv.Itoa(day)
	}

	if hour < 10 {
		strhour = strconv.Itoa(hour)
		strhour = "0" + strhour
	} else {
		strhour = strconv.Itoa(hour)
	}

	if min < 10 {
		strmin = strconv.Itoa(min)
		strmin = "0" + strmin
	} else {
		strmin = strconv.Itoa(min)
	}

	if sec < 10 {
		strsec = strconv.Itoa(sec)
		strsec = "0" + strsec
	} else {
		strsec = strconv.Itoa(sec)
	}

	origFilename := filepath.Base(path)
	// remove whitespace in previous files and replace with _
	origFilename = strings.Replace(origFilename, " ", "_", -1)
	newFilename := fmt.Sprintf("%s%s%s_%s%s%s_%s", stryear, strmonth, strday, strhour, strmin, strsec, origFilename)

	destpath := fmt.Sprintf("%s/%s-%s", destDir, stryear, strmonth)
	prevDir = destpath

	// create destination dir
	_, err = os.Stat(destpath)
	if err != nil {
		err = os.MkdirAll(destpath, 0755)
		if err != nil {
			log.Fatalf("Failed to make destination dir: %s due to error: %v", destpath, err)
		}
	}

	destFile := fmt.Sprintf("%s/%s", destpath, newFilename)
	df, err := os.Create(destFile)
	if err != nil {
		log.Fatalf("Failed to create destination file: %v", err)
	}
	defer df.Close()

	// seek back to the beginning for kicks in case we are not at offset 0
	f.Seek(0, 0)
	written, err := io.Copy(df, f)
	if err != nil {
		log.Fatalf("Failed to copy to destination file %s due to error: %v", df.Name(), err)
	}

	log.Printf("Copied %d bytes %s to %s", written, f.Name(), df.Name())
	return nil
}

func numMonth(month time.Month) string {
	switch month {
	case time.January:
		return "01"
	case time.February:
		return "02"
	case time.March:
		return "03"
	case time.April:
		return "04"
	case time.May:
		return "05"
	case time.June:
		return "06"
	case time.July:
		return "07"
	case time.August:
		return "08"
	case time.September:
		return "09"
	case time.October:
		return "10"
	case time.November:
		return "11"
	case time.December:
		return "12"
	}

	return ""
}

// copyMissingEXIF is called to copy .JPG files into the last known good dir
func copyMissingEXIF(sf *os.File) {

	sf.Seek(0, 0)

	sourceFilename := filepath.Base(sf.Name())

	df, err := os.Create(prevDir + "/" + sourceFilename)
	if err != nil {
		log.Fatalf("Failed to create destination file: %s due to error: %v", df.Name(), err)
	}
	defer df.Close()

	written, err := io.Copy(df, sf)
	if err != nil {
		log.Fatalf("Failed to copy to destination file %s due to error: %v", df.Name(), err)
	}

	log.Printf("Copied %d bytes %s to %s", written, sf.Name(), df.Name())

}

// copyInvalidExtension copies to prevDir all non JPG files
func copyInvalidExtension(sourcefile string) {

	sourceFilename := filepath.Base(sourcefile)

	df, err := os.Create(prevDir + "/" + sourceFilename)
	if err != nil {
		log.Fatalf("Failed to create destination file: %s due to error: %v", df.Name(), err)
	}
	defer df.Close()

	sf, err := os.Open(sourcefile)
	if err != nil {
		log.Fatalf("Failed to open sourcefile: %s due to error: %v", sourcefile, err)
	}
	defer sf.Close()

	written, err := io.Copy(df, sf)
	if err != nil {
		log.Fatalf("Failed to copy to destination file %s due to error: %v", df.Name(), err)
	}

	log.Printf("Copied %d bytes %s to %s", written, sf.Name(), df.Name())

}
