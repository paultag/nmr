package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Tempdir() (func(), error) {
	popdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	name, err := ioutil.TempDir("", "nmr.")
	if err != nil {
		return nil, err
	}
	err = os.Chdir(name)
	if err != nil {
		return nil, err
	}

	return func() {
		err := os.Chdir(popdir)
		if err != nil {
			fmt.Printf("Error during tmpdir cleanup!: %s", err)
		}
		err = os.RemoveAll(name)
		if err != nil {
			fmt.Printf("Error during tmpdir cleanup!: %s", err)
		}
	}, nil
}
