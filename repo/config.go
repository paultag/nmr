package repo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"pault.ag/go/debian/control"
	"pault.ag/go/resolver"
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
	Arches         []string
	Schroot        string
}

func (n NMRConfig) LoadIndex(dist string) (*resolver.Candidates, error) {
	can := resolver.Candidates{}
	dists, err := LoadDistributions(n.Basedir)
	if err != nil {
		return nil, err
	}

	repoDistConfig, err := dists.GetDistConfig(dist)

	d, err := n.GetDistConfig(dist)
	if err != nil {
		return nil, err
	}

	for _, arch := range repoDistConfig.Architectures {
		err := resolver.AppendBinaryIndex(
			&can,
			n.Global.PublicArchiveRoot,
			dist,
			"main",
			arch,
		)
		if err != nil {
			return nil, err
		}
	}

	for _, arch := range repoDistConfig.Architectures {
		err := resolver.AppendBinaryIndex(
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
	Basedir string

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
	ret := NMRConfig{
		Basedir: basedir,
	}

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
			Arches:         strings.Split(para.Values["Arches"], " "),
			Upstream: UpstreamLocation{
				Root: upstream[0],
				Dist: upstream[1], // FIXME
			},
			Schroot: para.Values["Schroot"],
		})
	}

	return &ret, nil
}
