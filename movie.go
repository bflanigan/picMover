package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type videoMetadata struct {
	camera       string
	year         string
	month        string
	day          string
	hour         string
	minute       string
	second       string
	origFilename string
	newFilename  string
}

func examineMov(path string) {

	v := videoMetadata{
		origFilename: filepath.Base(path),
		camera:       "unknown",
		year:         "unknown",
		month:        "unknown",
	}

	var destinationDir string
	if len(movieDir) == 0 {
		destinationDir = prevDir
	} else {
		destinationDir = movieDir
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open sourcefile: %s due to error: %v", path, err)
	}
	defer f.Close()

	extractMovStamp(path, &v)
	s := strings.Replace(v.origFilename, " ", "_", -1)

	var newFilename string
	if v.year != "unknown" { // if we were able to get some sort of timestamp, then include it in the new filename
		newFilename = fmt.Sprintf("%s%s%s_%s%s%s_%s", v.year, v.month, v.day, v.hour, v.minute, v.second, s)
	} else { // if we could not get a timestamp, then just reuse the original filename
		newFilename = s
	}

	destpath := fmt.Sprintf("%s/%s-%s", destinationDir, v.year, v.month)

	copyFile(f, newFilename, destpath+"/"+v.camera)

}

func extractMovStamp(path string, v *videoMetadata) {

	log.Printf("Working on: %s", path)

	cmd := exec.Command(mediainfo, path)
	var stdoutbuf bytes.Buffer
	var stderrbuf bytes.Buffer

	cmd.Stderr = &stderrbuf
	cmd.Stdout = &stdoutbuf

	err := cmd.Run()
	if err != nil {
		log.Printf("Failed running mediainfo on %s due to error: %v", path, err)
		return
	}

	var gotDate bool
	s := bufio.NewScanner(bytes.NewReader(stdoutbuf.Bytes()))

	/*
		com.apple.quicktime.make                 : Apple
		com.apple.quicktime.model                : iPhone 7
		Writing application                      : HandBrake 0.9.4 2009112300
		Writing application                      : GoPro
		Writing application                      : CanonMVI06
		Movie_More                               : EASTMAN KODAK COMPANY  KODAK CX6330 ZOOM DIGITAL CAMERA
	*/

	for s.Scan() {
		// 		Encoded date                             : UTC 2016-10-30 11:51:25
		if !gotDate && strings.HasPrefix(s.Text(), "Encoded date") {
			fields := strings.Fields(s.Text())
			ymd := fields[4]
			hms := fields[5]
			s1 := strings.Split(ymd, "-")
			v.year = s1[0]
			v.month = s1[1]
			v.day = s1[2]
			s2 := strings.Split(hms, ":")
			v.hour = s2[0]
			v.minute = s2[1]
			v.second = s2[2]
			gotDate = true
		}
		//Mastered date                            : 2008/02/23/ 22:56
		//Mastered date                            : 2009-08-14 23:25:21
		//Mastered date                            : SAT MAY 01 13:08:24 2010
		if !gotDate && strings.HasPrefix(s.Text(), "Mastered date") {
			fields := strings.Split(s.Text(), ":")
			stamp := strings.TrimPrefix(fields[1], " ")

			if strings.Contains(stamp, "/") {
				//2008/02/23/ 22:56
				s1 := strings.Split(stamp, "/")
				v.year = s1[0]
				v.month = s1[1]
				v.day = s1[2]
				hm := strings.Fields(s.Text())
				s2 := strings.Split(hm[4], ":")
				v.hour = s2[0]
				v.minute = s2[1]
				v.second = "00"
			}

			if strings.Contains(stamp, "-") {
				//2009-08-14 23:25:21
				s1 := strings.Split(stamp, "-")
				v.year = s1[0]
				v.month = s1[1]
				v.day = s1[2]
				hm := strings.Fields(s.Text())
				s2 := strings.Split(hm[4], ":")
				v.hour = s2[0]
				v.minute = s2[1]
				v.second = s2[2]
			} else {
				//SAT MAY 01 13:08:24 2010
				fields := strings.Fields(s.Text())
				hms := fields[6]
				v.year = fields[7]
				v.month = numMonthString(fields[4])
				v.day = fields[5]
				s2 := strings.Split(hms, ":")
				v.hour = s2[0]
				v.minute = s2[1]
				v.second = s2[2]
			}

			gotDate = true
		}

		// 		com.apple.quicktime.creationdate         : 2016-10-30T12:51:21+0100
		if !gotDate && strings.HasPrefix(s.Text(), "com.apple.quicktime.creationdate") {
			fields := strings.Fields(s.Text())
			a := fields[2]
			b := strings.Split(a, "T")
			c := strings.Split(b[0], "-")
			v.year = c[0]
			v.month = c[1]
			v.day = c[2]
			d := strings.Split(b[1], ":")
			v.hour = d[0]
			v.minute = d[1]
			e := strings.Split(d[2], "+")
			v.second = e[0]
			gotDate = true
		}

		if strings.HasPrefix(s.Text(), "com.apple.quicktime.model") || strings.HasPrefix(s.Text(), "Writing application") || strings.HasPrefix(s.Text(), "Movie_More") {
			fields := strings.Split(s.Text(), ":")
			v.camera = fields[1]
			v.camera = strings.TrimPrefix(v.camera, " ")
		}
	}

	if err := s.Err(); err != nil {
		log.Fatalf("Scanner encountered error parsing mediainfo output: %v", err)
	}

}
