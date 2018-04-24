package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func examineMov(o *object) {

	v := videoMetadata{
		origFilename: o.SourceName,
		camera:       "unknown-mov",
		year:         "0000",
		month:        "00",
		day:          "00",
		hour:         "00",
		minute:       "00",
		second:       "00",
	}

	var err error
	o.FH, err = os.Open(o.FullSourcePath)
	if err != nil {
		log.Fatalf("Failed to open sourcefile: %s due to error: %v", o.FullSourcePath, err)
	}
	defer o.FH.Close()

	extractMovStamp(o.FullSourcePath, &v)
	s := strings.Replace(v.origFilename, " ", "_", -1)

	var newFilename, destpath string
	if v.year == "0000" {
		if len(unknownDir) == 0 {
			destpath = fmt.Sprintf("%s/movies/%s", prevDir, v.camera)
		} else {
			destpath = fmt.Sprintf("%s/movies/%s", unknownDir, v.camera)
		}
		newFilename = s
	} else {
		destpath = fmt.Sprintf("%s/%s-%s/movies/%s", destDir, v.year, v.month, v.camera)
		newFilename = fmt.Sprintf("%s%s%s_%s%s%s_%s", v.year, v.month, v.day, v.hour, v.minute, v.second, s)
	}

	if noRenameDest {
		copyFileNew(o, v.origFilename, destpath)
	} else {
		copyFileNew(o, newFilename, destpath)
	}

}

func extractMovStamp(path string, v *videoMetadata) {

	if debug {
		log.Printf("Working on: %s", path)
	}

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
			} else if strings.Contains(stamp, "-") {
				//2009-08-14 23:25:21
				//2010-03- 5 02:06:13
				s1 := strings.Split(stamp, "-")
				v.year = strings.TrimSpace(s1[0])
				v.month = strings.TrimSpace(s1[1])
				v.day = strings.TrimSpace(s1[2])
				hm := strings.Fields(s.Text())
				// hm[len(hm)-1] gets us the last field in the []string
				s2 := strings.Split(hm[len(hm)-1], ":")
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
			// shove Apple in front of the i Devices to match what we're extracting on the pictures
			// all other vendors are out of luck
			if strings.HasPrefix(v.camera, "iP") {
				v.camera = "Apple-" + v.camera
			}
		}
	}

	if err := s.Err(); err != nil {
		log.Fatalf("Scanner encountered error parsing mediainfo output: %v", err)
	}

	// extract only the first field of the day value
	if len(v.day) > 0 {
		fields := strings.Fields(v.day)
		v.day = fields[0]

		if len(v.day) == 1 {
			v.day = "0" + v.day
		}
	}

}
