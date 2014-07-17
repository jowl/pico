package main

import (
	"log"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type TimeInfo struct {
	DateTime          time.Time
	DateTimeOriginal  time.Time
	DateTimeDigitized time.Time
	ModTime           time.Time
}

type Picture struct {
	Path      string
	Timestamp *time.Time
}

type PictureBuilder struct {
	Input           chan string
	Output          chan *Picture
	SelectTimestamp func(ts *TimeInfo) *time.Time
}

func (p PictureBuilder) Run() {
	for path := range p.Input {
		timeInfo := extractTimeInfo(path)
		p.Output <- &Picture{
			Path:      path,
			Timestamp: p.SelectTimestamp(timeInfo),
		}
	}
	close(p.Output)
}

func extractTimeInfo(path string) (timeInfo *TimeInfo) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	timeInfo = new(TimeInfo)
	if exifData, err := exif.Decode(f); err == nil {
		if exifDateTime, err := exifData.Get(exif.DateTime); err == nil {
			timeInfo.DateTime = *parseExifDateTime(exifDateTime.StringVal())
		}

		if exifDateTimeOriginal, err := exifData.Get(exif.DateTimeOriginal); err == nil {
			timeInfo.DateTimeOriginal = *parseExifDateTime(exifDateTimeOriginal.StringVal())
		}

		if exifDateTimeDigitized, err := exifData.Get(exif.DateTimeDigitized); err == nil {
			timeInfo.DateTimeDigitized = *parseExifDateTime(exifDateTimeDigitized.StringVal())
		}
	}

	if fileInfo, err := f.Stat(); err == nil {
		timeInfo.ModTime = fileInfo.ModTime()
	}

	return
}

const exifLayout = "2006:01:02 15:04:05" // Mon Jan 2 15:04:05 -0700 MST 2006
var localLoc, _ = time.LoadLocation("Local")

func parseExifDateTime(dateTime string) *time.Time {
	t, err := time.ParseInLocation(exifLayout, dateTime, localLoc)
	if err != nil {
		return nil
	} else {
		return &t
	}
}
