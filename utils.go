package main

import (
	"io"
	"log"
	"os"
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

func parseExtensions(s string) []string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.ToLower(s)
	ss := strings.Split(s, ",")

	return ss
}
