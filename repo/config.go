package repo

import (
	"bufio"
	"os"
	"strings"

	"pault.ag/go/debian/control"
)

type GlobalConfig struct {
	control.Paragraph

	PublicArchiveRoot string
}

type UpstreamLocation struct {
	Root string
	Dist string
}

type DistConfig struct {
	control.Paragraph

	Names    []string
	Upstream UpstreamLocation
	Schroot  string
}

type NMRConfig struct {
	Global GlobalConfig
	Blocks []DistConfig
}

func LoadConfig(basedir string) (*NMRConfig, error) {
	ret := NMRConfig{}

	file, err := os.Open(basedir + "/conf/nmr")
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	global, err := control.ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	ret.Global = GlobalConfig{
		PublicArchiveRoot: global.Values["PublicArchiveRoot"],
	}

	for {
		para, err := control.ParseParagraph(reader)
		if err != nil {
			return nil, err
		}
		if para == nil {
			break
		}
		upstream := strings.Split(para.Values["Upstream"], " ")

		ret.Blocks = append(ret.Blocks, DistConfig{
			Names: strings.Split(para.Values["Name"], " "),
			Upstream: UpstreamLocation{
				Root: upstream[0],
				Dist: upstream[1], // FIXME
			},
			Schroot: para.Values["Schroot"],
		})
	}

	return &ret, nil
}
