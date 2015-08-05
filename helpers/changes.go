package helpers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"os"

	"fmt"
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
func LogChangesFromDsc(logPath, dscPath, suite, arch string) string {
	return ""
}

// Take a .changes file, fake a new .changes, and append the build log to the
// existing .changes file.
func AppendLogToChanges(logPath, changesPath string) string {
	return ""
}
