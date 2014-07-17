package main

import (
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
		base := path.Base(picture.Path)
		println(path.Join(o.Root, dir, base))
	}
	o.done <- true
}

func (o *Organizer) Await() {
	<- o.done
}
