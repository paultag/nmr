package helpers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"os"

	"pault.ag/go/debian/control"
)

func GenerateLogChangesFromDSC(dscPath, architecture, distribution, logPath string) (string, error) {
	dsc, err := control.ParseDscFile(dscPath)
	if err != nil {
		return "", err
	}
	return GenerateLogChanges(*dsc, architecture, distribution, logPath)
}

func GenerateLogChanges(source control.DSC, architecture string, distribution string, logPath string) (string, error) {
	stat, err := os.Stat(logPath)
	if err != nil {
		return "", err
	}
	size := stat.Size()

	md5Hash, err := HashFile(logPath, md5.New())
	if err != nil {
		return "", err
	}

	sha1Hash, err := HashFile(logPath, sha1.New())
	if err != nil {
		return "", err
	}

	sha256Hash, err := HashFile(logPath, sha256.New())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`Format: 1.8
Date: Wed, 29 Apr 2015 21:29:13 -0400
Source: %s
Binary: %s
Architecture: %s
Version: %s
Distribution: %s
Urgency: low
Maintainer: Fake Maintainer <fake@maintainer.com>
Changed-By: Fake Maintainer <fake@maintainer.com>
Description:
 This package is a fake shim to upload Log files.
Changes:
 This package is a fake shim to upload Log files.
Checksums-Sha1:
 %s %d %s
Checksums-Sha256:
 %s %d %s
Files:
 %s %d log extra %s
	`,
		source.Values["Source"],
		source.Values["Binary"],
		architecture,
		source.Values["Version"],
		distribution,

		sha1Hash, size, logPath,
		sha256Hash, size, logPath,
		md5Hash, size, logPath,
	), nil
}
