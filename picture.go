package main

import (
	"fmt"
	"log"
	"os"
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

func examinePic(o *object) {

	p := pictureMetadata{
		origFilename: o.SourceName,
		make:         "unknown",
		model:        "unknown",
	}

	unknown := "unknown-pic"

	var err error
	o.FH, err = os.Open(o.FullSourcePath)
	if err != nil {
		log.Fatalf("Failed to open %s for reading due to error: %v", o.FullSourcePath, err)
	}
	defer o.FH.Close()

	x, err := exif.Decode(o.FH)
	if err != nil {
		if debug {
			log.Printf("Failed to decode exif on file: %s", o.FullSourcePath)
		}

		// if we have not set a directory into which to copy files of unknown metadata,
		// just copy it into the last known good directory
		if len(unknownDir) > 0 {
			copyFileNew(o, o.SourceName, unknownDir+"/"+unknown)
		} else {
			copyFileNew(o, o.SourceName, prevDir+"/"+unknown)
		}
		return
	}

	Make, err := x.Get("Make")
	if err != nil {
		if debug {
			log.Printf("Failed to extract make from metadata for file: %s", o.FullSourcePath)
		}
	} else {
		p.make = strings.Replace(Make.String(), `"`, "", -1)
	}

	Model, err := x.Get("Model")
	if err != nil {
		if debug {
			log.Printf("Failed to extract model from metadata for file: %s", o.FullSourcePath)
		}
	} else {
		p.model = strings.Replace(Model.String(), `"`, "", -1)
	}

	// if we failed to get Make and Model, then go for the lens
	if p.make == "unknown" {
		Make, err := x.Get("LensMake")
		if err != nil {
			if debug {
				log.Printf("Failed to extract LensMake from metadata for file: %s", o.FullSourcePath)
			}
		} else {
			p.make = strings.Replace(Make.String(), `"`, "", -1)
		}
	}

	if p.model == "unknown" {
		Model, err := x.Get("LensModel")
		if err != nil {
			if debug {
				log.Printf("Failed to extract model from metadata for file: %s", o.FullSourcePath)
			}
		} else {
			//LensModel: "iPhone 6 back camera 4.15mm f/2.2"
			lm := strings.Replace(Model.String(), `"`, "", -1)
			fields := strings.Fields(lm)
			if len(fields) > 1 {
				p.model = fields[0] + " " + fields[1]
			} else {
				p.model = "unknown"
			}
		}
	}

	camera := fmt.Sprintf("%s-%s", p.make, p.model)

	dateTaken, err := x.DateTime()
	if err != nil {
		if debug {
			log.Printf("Failed to extract time/date from metadata on file: %s", o.FullSourcePath)
		}
		if len(unknownDir) > 0 {
			copyFileNew(o, o.SourceName, unknownDir+"/"+camera)
		} else {
			copyFileNew(o, o.SourceName, prevDir+"/"+camera)
		}
		return
	}

	//extract the timestamp from the jpg
	p.year, p.month, p.day, p.hour, p.minute, p.second = exifDecode(dateTaken)

	// remove whitespace in source filename and replace with _
	s := strings.Replace(p.origFilename, " ", "_", -1)
	newFilename := fmt.Sprintf("%s%s%s_%s%s%s_%s", p.year, p.month, p.day, p.hour, p.minute, p.second, s)

	destpath := fmt.Sprintf("%s/%s-%s", destDir, p.year, p.month)
	// we got the metadata from this picture file, set the global variable prevDir
	// so the next file we come across we can put into this same directory
	prevDir = destpath

	if noRenameDest {
		copyFileNew(o, o.SourceName, destpath+"/"+camera)
	} else {
		copyFileNew(o, newFilename, destpath+"/"+camera)
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
