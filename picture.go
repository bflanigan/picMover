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

type pictureMetadata struct {
	make         string
	model        string
	year         string
	month        string
	day          string
	hour         string
	minute       string
	second       string
	origFilename string
	newFilename  string
}

func examinePic(path string) {

	p := pictureMetadata{
		origFilename: filepath.Base(path),
		make:         "unknown",
		model:        "unknown",
	}

	unknown := "unknown-pic"

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %s for reading due to error: %v", path, err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		if debug {
			log.Printf("Failed to decode exif on file: %s", path)
		}

		// if we have not set a directory into which to copy files of unknown metadata,
		// just copy it into the last known good directory
		if len(unknownDir) > 0 {
			copyFile(f, p.origFilename, unknownDir+"/"+unknown)
		} else {
			copyFile(f, p.origFilename, prevDir+"/"+unknown)
		}
		return
	}

	Make, err := x.Get("Make")
	if err != nil {
		if debug {
			log.Printf("Failed to extract make from metadata for file: %s", path)
		}
	} else {
		p.make = strings.Replace(Make.String(), `"`, "", -1)
	}

	Model, err := x.Get("Model")
	if err != nil {
		if debug {
			log.Printf("Failed to extract model from metadata for file: %s", path)
		}
	} else {
		p.model = strings.Replace(Model.String(), `"`, "", -1)
	}

	camera := fmt.Sprintf("%s-%s", p.make, p.model)

	dateTaken, err := x.DateTime()
	if err != nil {
		if debug {
			log.Printf("Failed to extract time/date from metadata on file: %s", path)
		}
		if len(unknownDir) > 0 {
			copyFile(f, p.origFilename, unknownDir+"/"+camera)
		} else {
			copyFile(f, p.origFilename, prevDir+"/"+camera)
		}
		return
	}

	//extract the timestamp from the jpg
	p.year, p.month, p.day, p.hour, p.minute, p.second = exifDecode(dateTaken)

	// remove whitespace in previous files and replace with _
	s := strings.Replace(p.origFilename, " ", "_", -1)
	newFilename := fmt.Sprintf("%s%s%s_%s%s%s_%s", p.year, p.month, p.day, p.hour, p.minute, p.second, s)

	destpath := fmt.Sprintf("%s/%s-%s", destDir, p.year, p.month)
	// we got the metadata from this picture file, set the global variable prevDir
	// so the next file we come across we can put into this same directory
	prevDir = destpath

	if noRename {
		copyFile(f, p.origFilename, destpath+"/"+camera)
	} else {
		copyFile(f, newFilename, destpath+"/"+camera)
	}

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
