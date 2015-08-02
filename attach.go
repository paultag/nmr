package main

import (
	"os"

	"pault.ag/go/nmr/helpers"
)

func main() {
	err := helpers.MergeLogChangesFromDSC(
		"dput-ng_1.9_all.changes",
		os.Args[1],
		"all",
		"unstable",
		os.Args[2],
	)
	if err != nil {
		panic(err)
	}
}
