package main

import (
	"log"

	"pault.ag/go/nmr/repo"
	"pault.ag/go/reprepro"
)

func main() {
	log.Printf("Loading local and remote indexes.\n")

	c, _ := repo.LoadConfig("/home/tag/tmp/repo")
	i, e := c.LoadIndex("unstable")

	if e != nil {
		panic(e)
	}

	log.Printf("%d binary package names found and loaded\n", len(*i))

	repRepo := reprepro.NewRepo("/home/tag/tmp/repo")
	needsBuild, err := repRepo.BuildNeeding("unstable", "any")

	if err != nil {
		panic(err)
	}

	log.Printf("Computing build candidates\n")

	for _, status := range repo.ComputeBuildStatus(*repRepo, *i, needsBuild) {
		log.Printf("%s - %s (%s)\n", status.Package.Location, status.Buildable, status.Why)
	}
}
