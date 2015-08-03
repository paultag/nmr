package main

import (
	"fmt"

	// "pault.ag/go/nmr/archive"
	"pault.ag/go/nmr/repo"
	// "pault.ag/go/reprepro"
)

func main() {
	c, _ := repo.LoadConfig("/home/tag/tmp/repo")
	i, e := c.LoadIndex("sid")
	fmt.Printf("%s\n", e)
	fmt.Printf("%s %s\n", len(*i), e)
	return

	// rRepo := reprepro.NewRepo("/home/tag/tmp/repo")
	// needsBuild, err := rRepo.BuildNeeding("unstable", "any")
	// if err != nil {
	// 	panic(err)
	// }

	// for _, status := range repo.ComputeBuildStatus(*rRepo, *cans, needsBuild) {
	// 	fmt.Printf("%s - %s (%s)\n", status.Package.Location, status.Buildable, status.Why)
	// }
}
