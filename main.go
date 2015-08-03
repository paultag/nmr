package main

import (
	"fmt"

	"pault.ag/go/nmr/archive"
	"pault.ag/go/nmr/repo"
	"pault.ag/go/reprepro"
)

func main() {
	c, _ := repo.LoadConfig("/home/tag/tmp/repo")
	i, e := c.LoadIndex("sid")
	fmt.Printf("%s %s\n", i, e)
	return

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

	rRepo := reprepro.NewRepo("/home/tag/tmp/repo")
	needsBuild, err := rRepo.BuildNeeding("unstable", "any")
	if err != nil {
		panic(err)
	}

	for _, status := range repo.ComputeBuildStatus(*rRepo, *cans, needsBuild) {
		fmt.Printf("%s - %s (%s)\n", status.Package.Location, status.Buildable, status.Why)
	}
}
