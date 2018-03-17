package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func copyFile(f *os.File, destFilename string, destPath string) {

	// create destination dir
	_, err := os.Stat(destPath)
	if err != nil {
		err = os.MkdirAll(destPath, 0755)
		if err != nil {
			log.Fatalf("Failed to make destination dir: %s due to error: %v", destPath, err)
		}
	}

	f.Seek(0, 0)

	// check for destination file already present, if so, rename the one we are trying to copy
	_, err = os.Stat(destPath + "/" + destFilename)
	if err == nil {
		epoch := strconv.FormatInt(time.Now().UnixNano(), 10)
		destFilename = epoch + "-" + destFilename
		log.Printf("Detected filename conflict - renaming destination file: %s", destPath+"/"+destFilename)
	}

	df, err := os.Create(destPath + "/" + destFilename)
	if err != nil {
		log.Fatalf("Failed to create destination file: %s due to error: %v", df.Name(), err)
	}
	defer df.Close()

	written, err := io.Copy(df, f)
	if err != nil {
		log.Fatalf("Failed to copy from %s to destination file %s due to error: %v", f.Name(), df.Name(), err)
	}

	log.Printf("Copied %d bytes %s to %s", written, f.Name(), df.Name())
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

func numMonthString(s string) string {
	switch s {
	case "JAN":
		return "01"
	case "FEB":
		return "02"
	case "MAR":
		return "03"
	case "APR":
		return "04"
	case "MAY":
		return "05"
	case "JUN":
		return "06"
	case "JUL":
		return "07"
	case "AUG":
		return "08"
	case "SEP":
		return "09"
	case "OCT":
		return "10"
	case "NOV":
		return "11"
	case "DEC":
		return "12"
	}

	return ""
}

func parseExtensions(s string) []string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.ToLower(s)
	ss := strings.Split(s, ",")

	return ss
}
