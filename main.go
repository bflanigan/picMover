package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {

	var startDir, destDir string
	flag.StringVar(&startDir, "startDir", "", "Start traversing from here")
	flag.StringVar(&destDir, "destDir", "", "Put renamed files here")
	flag.Parse()

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
		log.Printf("Got error entering walkFunc for file %s, error: %v\n", path, err)
		return nil
	}

	if info.IsDir() {
		log.Printf("Skipping dir: %s", path)
		return nil
	}

	if info.Size() < 8192 {
		log.Printf("Skipping small file: %s with size: %d", path, info.Size())
		return nil
	}

	fields := strings.Split(info.Name(), `.`)
	lastField := len(fields) - 1
	extension := fields[lastField]
	lowerExtension := strings.ToLower(extension)

	if lowerExtension != "jpg" {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %s for reading due to error: %v", path, err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		log.Printf("Failed to decode exif from file: %s due to error: %v", path, err)
	}

	dateTaken, err := x.DateTime()
	if err != nil {
		log.Printf("Failed to extract Date + Time pic was taken on file: %s due to error: %v", path, err)
		return nil
	}

	year, month, day := dateTaken.Date()

	stryear := strconv.Itoa(year)
	strmonth := numMonth(month)

	var strday string
	if day < 10 {
		strday = strconv.Itoa(day)
		strday = "0" + strday
	} else {
		strday = strconv.Itoa(day)
	}

	hour, min, sec := dateTaken.Clock()
	var strhour, strmin, strsec string

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

	stamp := fmt.Sprintf("%s%s%s_%s%s%s_%s", stryear, strmonth, strday, strhour, strmin, strsec, origFilename)

	log.Printf("File: %s NewName: %v, ExtractedData: %s", path, stamp, dateTaken.String())
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
