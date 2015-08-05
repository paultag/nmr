package helpers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"pault.ag/go/debian/control"
)

func FakeChanges(
	date,
	source,
	binary,
	arch,
	version,
	distribution,
	urgency string,
	files []string,
) (string, error) {

	sha1FileHashes := ""
	sha256FileHashes := ""
	md5FileHashes := ""

	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			return "", err
		}
		size := stat.Size()

		md5Hash, err := HashFile(file, md5.New())
		if err != nil {
			return "", err
		}

		md5FileHashes += fmt.Sprintf(
			"\n %s %d fake extra %s",
			md5Hash,
			size,
			file,
		)

		sha1Hash, err := HashFile(file, sha1.New())
		if err != nil {
			return "", err
		}

		sha1FileHashes += fmt.Sprintf(
			"\n %s %d %s",
			sha1Hash,
			size,
			file,
		)

		sha256Hash, err := HashFile(file, sha256.New())
		if err != nil {
			return "", err
		}

		sha256FileHashes += fmt.Sprintf(
			"\n %s %d %s",
			sha256Hash,
			size,
			file,
		)
	}

	return fmt.Sprintf(`Format: 1.8
Date: %s
Source: %s
Binary: %s
Architecture: %s
Version: %s
Distribution: %s
Urgency: %s
Maintainer: Fake Maintainer <fake-maintainer@example.com>
Changed-By: Fake Maintainer <fake-maintainer@example.com>
Description:
 fake changes file to trick things into thinking something
 actually happened, even though it didn't.
Changes:
 fake changes file to trick things into thinking something
 actually happened, even though it didn't.
Checksums-Sha1: %s
Checcksums-Sha256: %s
Files: %s
`,
		date,
		source,
		binary,
		arch,
		version,
		distribution,
		urgency,
		// checksums

		sha1FileHashes,
		sha256FileHashes,
		md5FileHashes,
	), nil
}

// Take a DSC, and a log file, and create a fake .changes to upload
// the log to the archive. This relies on a reprepro extension.
func LogChangesFromDsc(logPath, dscPath, suite, arch string) (string, error) {
	dsc, err := control.ParseDscFile(dscPath)
	if err != nil {
		return "", nil
	}

	return FakeChanges(
		"Fri, 31 Jul 2015 12:53:50 -0400",
		dsc.Source,
		strings.Join(dsc.Binaries, " "),
		arch,
		dsc.Version.String(),
		suite,
		"low",
		[]string{logPath},
	)
}

func LogChangesFromChanges(logPath, changesPath, arch string) (string, error) {
	changes, err := control.ParseChangesFile(changesPath)
	if err != nil {
		return "", nil
	}

	return FakeChanges(
		"Fri, 31 Jul 2015 12:53:50 -0400",
		changes.Source,
		strings.Join(changes.Binaries, " "),
		arch,
		changes.Version.String(),
		changes.Distribution,
		changes.Urgency,
		[]string{logPath},
	)
}

func AppendLogToChanges(logPath, changesPath, arch string) error {
	changes, err := LogChangesFromChanges(logPath, changesPath, arch)
	if err != nil {
		return err
	}
	f, err := ioutil.TempFile("", "nmr.")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	_, err = f.Write([]byte(changes))
	if err != nil {
		return err
	}

	changes, err = MergeChanges(changesPath, f.Name())
	if err != nil {
		return err
	}

	fd, err := os.Create(changesPath)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write([]byte(changes))
	if err != nil {
		return err
	}

	return nil
}

func MergeChanges(changes ...string) (string, error) {
	bytes, err := exec.Command("mergechanges", changes...).Output()
	return string(bytes), err
}
