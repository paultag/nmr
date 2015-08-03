package archive

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"pault.ag/go/nmr/candidate"
)

func GetBinaryIndex(mirror, suite, component, arch string) (*candidate.Canidates, error) {
	resp, err := http.Get(fmt.Sprintf(
		"%s/dists/%s/%s/binary-%s/Packages.gz",
		mirror, suite, component, arch,
	)) // contains arch:all in amd64, etc
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return candidate.ReadFromBinaryIndex(reader)
}
