package main

import (
	"fmt"
	"os"
	"path"
)

type Organizer struct {
	input         chan *Picture
	root          string
	dryRun        bool
	done          chan bool
	timestampPath func(*Picture) (string, error)
}

func (o Organizer) Run() {
	for picture := range o.input {
		dir, err := o.timestampPath(picture)
		if err != nil {
			LogWarningf("Won't move %v: %v", picture.Path, err)
			continue
		}
		fname := path.Base(picture.Path)
		fullPath := path.Join(o.root, dir, fname)
		move(picture.Path, fullPath, o.dryRun)
	}
	o.done <- true
}

func (o Organizer) Await() {
	<-o.done
}

func move(source, target string, dryRun bool) {
	var err error
	if _, statErr := os.Stat(target); os.IsNotExist(statErr) {
		if !dryRun {
			if err = os.MkdirAll(path.Dir(target), 0755); err == nil {
				if err = os.Rename(source, target); err == nil {
					LogInfof("Moved %v to %v", source, target)
				}
			}
		} else {
			fmt.Printf("%v would be moved to %v\n", source, target)
		}
	} else if (err == nil) {
		err = os.ErrExist
	}


	if err != nil {
		LogErrorf("Couldn't move %v: %v", source, err)
	}
}
