package main

import (
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type TimeInfo struct {
	DateTime          *time.Time
	DateTimeOriginal  *time.Time
	DateTimeDigitized *time.Time
	ModTime           *time.Time
}

type Picture struct {
	Path      string
	Timestamp *time.Time
}

type PictureBuilder struct {
	input           chan string
	output          chan *Picture
	selectTimestamp func(ts *TimeInfo) *time.Time
}

func (p PictureBuilder) Run() {
	for path := range p.input {
		timeInfo := extractTimeInfo(path)
		if timeInfo == nil {
			continue
		}
		p.output <- &Picture{
			Path:      path,
			Timestamp: p.selectTimestamp(timeInfo),
		}
	}
	close(p.output)
}

func extractTimeInfo(path string) (timeInfo *TimeInfo) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		LogErrorf("Couldn't open file to extract time info: %v", err)
		return nil
	}
	timeInfo = new(TimeInfo)
	if exifData, err := exif.Decode(f); err == nil {
		if exifDateTime, err := exifData.Get(exif.DateTime); err == nil {
			timeInfo.DateTime = parseExifDateTime(exifDateTime.StringVal())
		} else {
			LogWarningf("Couldn't parse EXIF DateTime: %v", err)
		}

		if exifDateTimeOriginal, err := exifData.Get(exif.DateTimeOriginal); err == nil {
			timeInfo.DateTimeOriginal = parseExifDateTime(exifDateTimeOriginal.StringVal())
		} else {
			LogWarningf("Couldn't parse EXIF DateTimeOriginal: %v", err)
		}

		if exifDateTimeDigitized, err := exifData.Get(exif.DateTimeDigitized); err == nil {
			timeInfo.DateTimeDigitized = parseExifDateTime(exifDateTimeDigitized.StringVal())
		} else {
			LogWarningf("Couldn't parse EXIF DateTimeDigitized: %v", err)
		}
	} else {
		LogWarningf("Couldn't parse EXIF data: %v", err)
	}

	if fileInfo, err := f.Stat(); err == nil {
		modTime := fileInfo.ModTime()
		timeInfo.ModTime = &modTime
	}

	return
}

const exifLayout = "2006:01:02 15:04:05" // Mon Jan 2 15:04:05 -0700 MST 2006

func parseExifDateTime(dateTime string) *time.Time {
	t, err := time.ParseInLocation(exifLayout, dateTime, time.Local)
	if err != nil {
		return nil
	} else {
		return &t
	}
}
