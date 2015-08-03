package main

import (
	"fmt"

	"pault.ag/go/nmr/archive"
	"pault.ag/go/nmr/build"
	"pault.ag/go/reprepro"
)

func main() {
	cans, err := archive.GetBinaryIndex(
		"http://http.debian.net/debian",
		"unstable",
		"main",
		"amd64",
	)
	if err != nil {
		panic(err)
	}

	err = archive.AppendBinaryIndex(
		cans,
		"https://pault.ag/debian",
		"wicked",
		"main",
		"amd64",
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", (*cans)["hairycandy-serving-suggestions"])
	fmt.Printf("%s\n", (*cans)["fluxbox"])
	return

	repo := reprepro.NewRepo("/home/tag/tmp/repo")
	needsBuild, err := repo.BuildNeeding("unstable", "any")
	if err != nil {
		panic(err)
	}

	for _, status := range build.ComputeBuildStatus(*repo, *cans, needsBuild) {
		fmt.Printf("%s - %s (%s)\n", status.Package.Location, status.Buildable, status.Why)
	}
}
