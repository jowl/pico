package main

import (
	"bufio"
	"os"
	"time"

	"github.com/codegangsta/cli"
)

type Pico struct {
	SelectTimestamp func(*TimeInfo) *time.Time
	TimestampPath   func(*Picture) string
	app             *cli.App
}

func New() *Pico {
	pico := &Pico{
		SelectTimestamp: defaultSelectTimestamp,
		TimestampPath:   defaultTimestampPath,
	}

	app := cli.NewApp()
	app.Name = "pico"
	app.Usage = "picture organizer"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "target-dir,d", Value: ".", Usage: "target root directory"},
		cli.BoolFlag{Name: "dry-run,n", Usage: "show which files would have been moved"},
	}
	app.Action = func(c *cli.Context) {
		pico.run(c)
	}

	pico.app = app
	return pico
}

func (p Pico) Run(args []string) {
	p.app.Run(args)
}

func (p Pico) run(c *cli.Context) {
	input := make(chan string)
	pictures := make(chan *Picture)

	pictureBuilder := &PictureBuilder{
		input:           input,
		output:          pictures,
		selectTimestamp: p.SelectTimestamp,
	}

	organizer := &Organizer{
		input:         pictures,
		root:          c.String("target-dir"),
		dryRun:        c.Bool("dry-run"),
		done:          make(chan bool),
		timestampPath: p.TimestampPath,
	}

	go pictureBuilder.Run()
	go organizer.Run()

	if len(c.Args()) > 0 {
		for _, path := range c.Args() {
			input <- path
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input <- scanner.Text()
		}
	}
	close(input)

	organizer.Await()
}

func defaultSelectTimestamp(ti *TimeInfo) *time.Time {
	t := ti.DateTimeOriginal
	if t == nil || t.Before(time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC)) {
		return nil
	}
	return t
}

func defaultTimestampPath(p *Picture) string {
	if p.Timestamp == nil {
		return "unknown"
	}
	return p.Timestamp.Format("2006/2006-01") // Mon Jan 2 15:04:05 -0700 MST 2006
}
