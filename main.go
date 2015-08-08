package main

import (
	"fmt"
	"os"

	"pault.ag/go/debian/control"
	// "pault.ag/go/fancytext"
	"pault.ag/go/nmr/helpers"
	"pault.ag/go/reprepro"
)

func main() {
	arch := os.Args[1]
	repoRoot := os.Args[2]
	params := os.Args[3:]

	repreproRepo := reprepro.NewRepo(repoRoot)
	log, err := repreproRepo.ParseLogEntry(params)
	if err != nil {
		panic(err)
	}

	arches := log.Changes.Architectures
	// insane hack here. ignore me.
	fmt.Printf("Got: %s %s\n", log.Action, log.Changes.Architectures)
	if log.Action != "accepted" || arches[0].CPU != "source" || len(arches) != 1 {
		fmt.Printf("Ignoring: %s %s\n", log.Action, log.Changes.Architectures)
		return
	}

	// Right, let's run a build.
	dsc, err := log.Changes.GetDSC()
	if err != nil {
		panic(err)
	}

	complete, err := Tempdir()
	if err != nil {
		panic(err)
	}
	defer complete()

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

	err = changes.Copy(incomingLocation)
	if err != nil {
		return ftbfs, err
	}

	return ftbfs, nil
}
