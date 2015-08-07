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

	complete, err := Tempdir()
	if err != nil {
		panic(err)
	}
	defer complete()

	params := os.Args[1:]
	log, err := ParseLine(repoRoot, params)
	if err != nil {
		panic(err)
	}

	arches := log.Changes.Architectures
	// insane hack here. ignore me.
	if log.Action != "accepted" && len(arches) == 1 && arches[0].CPU == "source" {
		fmt.Printf("Ignoring: %s %s\n", log.Action, log.Changes.Architectures)
		return
	}

	// Right, let's run a build.
	dsc, err := log.Changes.GetDSC()
	if err != nil {
		panic(err)
	}

	ftbfs, err := BuildPackage(dsc.Filename, arch, log.Suite, repoRoot, true)
	fmt.Printf("FTBFS: %s", ftbfs)
	fmt.Printf("Error: %s", err)
}

func BuildPackage(dscFile, arch, suite, repoRoot string, verbose bool) (bool, error) {
	// done := fancytext.FormatSpinner(fmt.Sprintf("%%s  -  building %s", dscFile))
	// defer done()

	incomingLocation, err := GetIncoming(repoRoot, suite)
	if err != nil {
		return false, err
	}

	dsc, err := control.ParseDscFile(dscFile)
	if err != nil {
		return false, err
	}

	source := dsc.Source
	version := dsc.Version

	cmd, err := SbuildCommand(suite, suite, arch, dscFile, repoRoot, verbose)
	if err != nil {
		return false, err
	}

	err = cmd.Run()
	ftbfs := err != nil

	changesFile := helpers.Filename(source, version, arch, "changes")
	logPath := helpers.Filename(source, version, arch, "build")

	if ftbfs {
		return true, nil

		// fmt.Printf(" FTBFS!\n")
		// return
		// changes, err := helpers.LogChangesFromDsc(logPath, dscFile, suite, arch)
		// if err != nil {
		// 	panic(err)
		// }
		// fd, err := os.Create(changesFile)
		// if err != nil {
		// 	panic(err)
		// }
		// defer fd.Close()
		// _, err = fd.Write([]byte(changes))
		// if err != nil {
		// 	panic(err)
		// }
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
		return ftbfs, err
	}

	err = changes.Move(incomingLocation)
	if err != nil {
		return ftbfs, err
	}

	return ftbfs, nil
}
