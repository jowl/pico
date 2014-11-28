package main

import (
	"bytes"
	"fmt"
	"io"
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

func hasSameContent(path1, path2 string) (bool, error) {
	b1 := make([]byte, 10)
	b2 := make([]byte, 10)
	f1, _ := os.Open(path1)
	f2, _ := os.Open(path2)
	var e1, e2 error
	for {
		_, e1 = f1.Read(b1)
		_, e2 = f2.Read(b2)
		if e1 != e2 || !bytes.Equal(b1, b2) {
			if e1 != nil {
				return false, e1
			} else {
				return false, e2
			}
		}
		if e1 == io.EOF {
			break
		}
	}
	return true, nil;
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
		if same, sameErr := hasSameContent(source, target); sameErr == nil && same {
			err = fmt.Errorf("file already exists, but has the same content as %v", source)
		} else {
			err = os.ErrExist
		}
	}

	if err != nil {
		LogErrorf("Couldn't move %v to %v: %v", source, target, err)
	}
}
