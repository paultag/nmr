package repo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"pault.ag/go/debian/control"
	"pault.ag/go/nmr/archive"
	"pault.ag/go/nmr/candidate"
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

	Names          []string
	Upstream       UpstreamLocation
	UpstreamArches []string
	Schroot        string
}

func (d DistConfig) LoadIndex() (*candidate.Canidates, error) {
	can := candidate.Canidates{}
	for _, arch := range d.UpstreamArches {
		err := archive.AppendBinaryIndex(
			&can,
			d.Upstream.Root,
			d.Upstream.Dist,
			"main",
			arch,
		)
		if err != nil {
			return nil, err
		}
	}
	return &can, nil
}

type NMRConfig struct {
	Global GlobalConfig
	Blocks []DistConfig
}

func (nmr NMRConfig) GetDistConfig(name string) (*DistConfig, error) {
	for _, block := range nmr.Blocks {
		for _, dname := range block.Names {
			if dname == name {
				return &block, nil
			}
		}
	}
	return nil, fmt.Errorf("No such name: %s", name)
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
			Names:          strings.Split(para.Values["Name"], " "),
			UpstreamArches: strings.Split(para.Values["UpstreamArches"], " "),
			Upstream: UpstreamLocation{
				Root: upstream[0],
				Dist: upstream[1], // FIXME
			},
			Schroot: para.Values["Schroot"],
		})
	}

	return &ret, nil
}
