package main

import (
	"fmt"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/version"
)

/*
accepted sid liblicense 0.8.1-3 /home/tag/tmp/repo/tmp/liblicense_0.8.1-3_source.changes pool/main/libl/liblicense/liblicense_0.8.1-3_source.changes
*/

type Log struct {
	Action  string
	Suite   string
	Source  string
	Version version.Version
	Changes control.Changes
}

func ParseLine(root string, params []string) (*Log, error) {
	if len(params) != 6 {
		return nil, fmt.Errorf("Unknown input string format")
	}

	version, err := version.Parse(params[3])
	if err != nil {
		return nil, err
	}

	changes, err := control.ParseChangesFile(root + "/" + params[5])
	if err != nil {
		return nil, err
	}

	return &Log{
		Action:  params[0],
		Suite:   params[1],
		Source:  params[2],
		Version: version,
		Changes: *changes,
	}, nil
}
