package main

import (
	"bufio"
	"os"
	"time"

	"github.com/codegangsta/cli"
)

func Run(args []string) {
	app := cli.NewApp()
	app.Name = "pico"
	app.Usage = "picture organizer"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "target-dir,d", Value: ".", Usage: "target root directory"},
		cli.BoolFlag{Name: "dry-run,n", Usage: "show which files would have been moved"},
	}
	app.Action = runPico
	app.Run(args)
}

func runPico(c *cli.Context) {
	input := make(chan string)
	pictures := make(chan *Picture)

	pictureBuilder := &PictureBuilder{
		Input:  input,
		Output: pictures,
		SelectTimestamp: func(ti *TimeInfo) *time.Time {
			t := ti.DateTimeOriginal
			if t.Before(time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC)) {
				return nil
			}
			return &t
		},
	}

	organizer := &Organizer{
		Input:  pictures,
		Root:   c.String("target-dir"),
		DryRun: c.Bool("dry-run"),
		done:   make(chan bool),
		TimestampPath: func(t *time.Time) string {
			if t == nil {
				return "unknown"
			}
			return t.Format("2006/2006-01") // Mon Jan 2 15:04:05 -0700 MST 2006
		},
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
