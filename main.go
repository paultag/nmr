package main

import (
	"fmt"
	"os"

	"pault.ag/go/debian/control"
	// "pault.ag/go/fancytext"
	"pault.ag/go/nmr/helpers"
)

func main() {
	arch := "amd64"
	repoRoot := "/home/tag/tmp/repo/"

	params := os.Args[1:]
	log, err := ParseLine(repoRoot, params)
	if err != nil {
		panic(err)
	}

	if log.Action != "accepted" {
		return
	}

	// Right, let's run a build.
	dsc, err := log.Changes.GetDSC()
	if err != nil {
		panic(err)
	}

	BuildPackage(dsc.Filename, arch, log.Suite, repoRoot, true)
}

func BuildPackage(dscFile, arch, suite, repoRoot string, verbose bool) {
	// done := fancytext.FormatSpinner(fmt.Sprintf("%%s  -  building %s", dscFile))
	// defer done()

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

	cmd, err := SbuildCommand(suite, suite, arch, dscFile, repoRoot, verbose)
	if err != nil {
		panic(err)
	}

	err = cmd.Run()
	ftbfs := err != nil

	changesFile := helpers.Filename(source, version, arch, "changes")
	logPath := helpers.Filename(source, version, arch, "build")

	if ftbfs {
		fmt.Printf(" FTBFS!\n")
		return
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

	if IsArchAllArch(repoRoot, arch) && dsc.HasArchAll() {
		archAllLogPath := helpers.Filename(source, version, "all", "build")
		Copy(logPath, archAllLogPath)
		helpers.AppendLogToChanges(archAllLogPath, changesFile, "all")
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
