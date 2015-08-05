package helpers

import (
	"fmt"

	"pault.ag/go/debian/version"
)

func Filename(source string, v version.Version, arch, flavor string) string {
	// file paths don't have the epoch in them.
	version := ""
	if v.IsNative() {
		version = v.Version
	} else {
		version = fmt.Sprintf("%s-%s", v.Version, v.Revision)
	}
	return fmt.Sprintf(
		"%s_%s_%s.%s",
		source, version, arch, flavor,
	)

}
