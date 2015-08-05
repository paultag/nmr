package main

import (
	"fmt"
	"os"

	"pault.ag/go/debian/control"
	"pault.ag/go/nmr/helpers"
)

func main() {
	arch := "amd64"
	dscFile := "/home/tag/tmp/repo/pool/main/f/fbautostart/fbautostart_2.718281828-4.dsc"
	suite := "unstable"
	repoRoot := "/home/tag/tmp/repo/"
	verbose := true

	incomingLocation, err := GetIncoming(repoRoot, suite)
	if err != nil {
		panic(err)
	}

	dsc, err := control.ParseDscFile(dscFile)
	if err != nil {
		panic(err)
	}

	source := dsc.Source
	version := dsc.Version

	fmt.Printf(`

		░█▀█░█▀█░█░█░█▀▀░█▀▄░░░█▄█░█▀█░█░░░█▀▀░░░█▀▄░█▀█░▀█▀
		░█░█░█▀█░█▀▄░█▀▀░█░█░░░█░█░█░█░█░░░█▀▀░░░█▀▄░█▀█░░█░
		░▀░▀░▀░▀░▀░▀░▀▀▀░▀▀░░░░▀░▀░▀▀▀░▀▀▀░▀▀▀░░░▀░▀░▀░▀░░▀░

		Source: %s
		Version: %s
		Arch: %s
		Suite: %s
		Chroot: %s
		Repo: %s

`, source, version, arch, suite, suite, repoRoot)

	cmd, err := SbuildCommand(suite, suite, arch, dscFile, repoRoot, verbose)
	if err != nil {
		panic(err)
	}

	err = cmd.Run()
	ftbfs := err != nil

	fmt.Printf("%s %s\n", err, ftbfs)
	changesFile := helpers.Filename(source, version, arch, "changes")
	logPath := helpers.Filename(source, version, arch, "build")

	if ftbfs {
		changes, err := helpers.LogChangesFromDsc(logPath, dscFile, suite, arch)
		if err != nil {
			panic(err)
		}
		fd, err := os.Create(changesFile)
		if err != nil {
			panic(err)
		}
		defer fd.Close()
		_, err = fd.Write([]byte(changes))
		if err != nil {
			panic(err)
		}
	} else {
		helpers.AppendLogToChanges(logPath, changesFile, arch)
	}

	changes, err := control.ParseChangesFile(changesFile)
	if err != nil {
		panic(err)
	}

	err = changes.Move(incomingLocation)
	if err != nil {
		panic(err)
	}
}
