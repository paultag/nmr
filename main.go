package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/inotify"
	"pault.ag/go/debian/control"
)

func main() {
	watchDirectory(".")
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
		if ev.Mask != 0x08 { // IN_CLOSE_WRITE
			continue
		}
		if !strings.HasSuffix(ev.Name, ".changes") {
			continue
		}
		f, err := os.Open(ev.Name)
		if err != nil {
			log.Printf("%s\n", err)
			continue
		}
		para, err := control.ParseChanges(bufio.NewReader(f))
		if err != nil {
			log.Printf("%s\n", err)
			continue
		}
		log.Printf("%s\n", para.Distribution)
	}

	return nil
}
