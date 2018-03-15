package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func examinePic(path string) {

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %s for reading due to error: %v", path, err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		if len(unknownDir) == 0 {
			copyFile(f, filepath.Base(path), prevDir)
		} else {
			copyFile(f, filepath.Base(path), unknownDir)
		}
		return
	}

	dateTaken, err := x.DateTime()
	if err != nil {
		if len(unknownDir) == 0 {
			copyFile(f, filepath.Base(path), prevDir)
		} else {
			copyFile(f, filepath.Base(path), unknownDir)
		}
		return
	}

	//extract the timestamp from the jpg
	stryear, strmonth, strday, strhour, strmin, strsec := exifDecode(dateTaken)

	// remove whitespace in previous files and replace with _
	origFilename := filepath.Base(path)
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

	copyFile(f, newFilename, destpath)

}

func exifDecode(dateTaken time.Time) (string, string, string, string, string, string) {
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

	return stryear, strmonth, strday, strhour, strmin, strsec
}
