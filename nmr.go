package main

import (
	"os"

	"pault.ag/go/nmr/helpers"
)

func main() {
	err := helpers.MergeLogChangesFromDSC(
		os.Args[1],
		os.Args[2],
		os.Args[3],
		os.Args[4],
		os.Args[5],
	)
	if err != nil {
		panic(err)
	}
}
