package main

import (
	"bufio"
	"errors"
	"os"
	"time"

	"github.com/codegangsta/cli"
)

type Pico struct {
	SelectTimestamp func(*TimeInfo) *time.Time
	TimestampPath   func(*Picture) (string, error)
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
		cli.BoolFlag{Name: "dry-run,n", Usage: "show which files would have been moved"},
	}
	app.Commands = []cli.Command{*pico.newOrganizeCommand()}
	pico.app = app
	return pico
}

func (p Pico) Run(args []string) {
	p.app.Run(args)
}

func (p Pico) newOrganizeCommand() *cli.Command {
	return &cli.Command{
		Name:      "organize",
		ShortName: "o",
		Usage:     "organize pictures from stdin/command-line arguments",
		Action: func(c *cli.Context) {
			p.organize(c)
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "target-dir,d", Value: ".", Usage: "target root directory"},
		},
	}
}

func (p Pico) organize(c *cli.Context) {
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
		dryRun:        c.GlobalBool("dry-run"),
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
	beginningOfTime := time.Date(1,1,1,0,0,0,0,time.UTC)
	if t := ti.DateTimeOriginal; t != nil && t.After(beginningOfTime) {
		return t
	} else if t := ti.DateTimeDigitized; t != nil && t.After(beginningOfTime) {
		return t
	} else if t := ti.DateTime; t != nil && t.After(beginningOfTime) {
		return t
	} else {
		return nil
	}
}

func defaultTimestampPath(p *Picture) (string, error) {
	if p.Timestamp == nil {
		return "", errors.New("Timestamp is nil")
	}
	return p.Timestamp.Format("2006/2006-01"), nil // Mon Jan 2 15:04:05 -0700 MST 2006
}
