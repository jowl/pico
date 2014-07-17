package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type Organizer struct {
	Input         chan *Picture
	Root          string
	DryRun        bool
	done          chan bool
	TimestampPath func(t *time.Time) string
}

func (o *Organizer) Run() {
	for picture := range o.Input {
		dir := o.TimestampPath(picture.Timestamp)
		fname := path.Base(picture.Path)
		fullPath := path.Join(o.Root, dir, fname)
		if o.DryRun {
			printMove(picture.Path, fullPath)
		} else {
			doMove(picture.Path, fullPath)
		}
	}
	o.done <- true
}

func (o *Organizer) Await() {
	<-o.done
}

func printMove(source, target string) {
	fmt.Printf("%v would be moved to %v\n", source, target)
}

func doMove(source, target string) {
	var err error
	if err = os.MkdirAll(path.Dir(target), 0755); err == nil {
		err = os.Rename(source, target)
	}
	if err != nil {
		log.Printf("Error when moving %v: %v", source, err)
	}
}
