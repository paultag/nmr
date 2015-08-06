package main

import (
	"fmt"
	"os"

	"pault.ag/go/debian/control"
	"pault.ag/go/fancytext"
	"pault.ag/go/nmr/helpers"
)

func main() {
	arch := "amd64"
	suite := "unstable"
	repoRoot := "/home/tag/tmp/repo/"

	fmt.Printf(`

		░█▀█░█▀█░█░█░█▀▀░█▀▄░░░█▄█░█▀█░█░░░█▀▀░░░█▀▄░█▀█░▀█▀
		░█░█░█▀█░█▀▄░█▀▀░█░█░░░█░█░█░█░█░░░█▀▀░░░█▀▄░█▀█░░█░
		░▀░▀░▀░▀░▀░▀░▀▀▀░▀▀░░░░▀░▀░▀▀▀░▀▀▀░▀▀▀░░░▀░▀░▀░▀░░▀░

				Arch:    %s
				Suite:   %s
				Chroot:  %s
				Repo:    %s

`, arch, suite, suite, repoRoot)

	fmt.Printf("Loading package indexes and computing build needed...\n")

	needsBuild := GetBuildNeeding(repoRoot, suite, arch)
	if IsArchAllArch(repoRoot, arch) {
		fmt.Printf("Also getting needs build packages for all\n")
		needsBuild = append(needsBuild, GetBuildNeeding(repoRoot, suite, "all")...)
	}

	fmt.Printf("  %d packages need build\n", len(needsBuild))
	fmt.Printf("\n\n")

	for _, pkg := range needsBuild {
		dscPath := repoRoot + "/" + pkg.Package.Location

		if pkg.Buildable {
			fmt.Printf("Building: %s", dscPath)
			BuildPackage(dscPath, arch, suite, repoRoot, false)
			fmt.Printf("      built %s!\n", dscPath)
		} else {
			fmt.Printf("Package %s unbuildable currently", pkg.Package.Location)
		}
	}
}

func BuildPackage(dscFile, arch, suite, repoRoot string, verbose bool) {
	done := fancytext.FormatSpinner(fmt.Sprintf("%%s  -  building %s", dscFile))
	defer done()

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
