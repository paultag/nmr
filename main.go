package main

import (
	"log"
	"strings"

	"golang.org/x/exp/inotify"
	"pault.ag/go/debian/control"
)

func main() {
	watchDirectory(".")
}

func processChanges(changes control.Changes) {
	dsc, err := changes.GetDSC()
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("%s\n", dsc.Filename)

	err = changes.Move("/home/tag/tmp/x/")
	if err != nil {
		log.Printf("%s\n", err)
	}
}

func watchDirectory(path string) error {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Watch(path)
	if err != nil {
		return err
	}

	for ev := range watcher.Event {
		if ev.Mask != 0x08 || !strings.HasSuffix(ev.Name, ".changes") {
			// 0x08 -> IN_CLOSE_WRITE
			// and wait for the .changes file
			continue
		}
		para, err := control.ParseChangesFile(ev.Name)
		if err != nil {
			log.Printf("%s\n", err)
			continue
		}
		go processChanges(*para)
	}

	return nil
}
