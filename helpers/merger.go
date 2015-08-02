package helpers

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func MergeChanges(changesOne, changesTwo string) error {
	cmd := exec.Command("mergechanges", changesOne, changesTwo)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fd, err := os.Create(changesOne)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func MergeLogChangesFromDSC(targetChanges, dscPath, architecture, distribution, logPath string) error {
	out, err := GenerateLogChangesFromDSC(dscPath, architecture, distribution, logPath)
	if err != nil {
		return err
	}

	fd, err := ioutil.TempFile(".", "nmr")
	fd.Write([]byte(out))

	if err != nil {
		return err
	}

	return MergeChanges(targetChanges, fd.Name())
}
