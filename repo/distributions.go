package repo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"pault.ag/go/debian/control"
)

type Distribution struct {
	control.Paragraph

	Codename      string
	Suite         string
	Components    []string
	Architectures []string
	Tracking      []string
	Source        bool
}

type Distributions struct {
	Blocks []Distribution
}

func (d Distributions) GetDistConfig(name string) (*Distribution, error) {
	for _, block := range d.Blocks {
		if block.Codename == name || block.Suite == name {
			return &block, nil
		}
	}
	return nil, fmt.Errorf("No such name: %s", name)
}

func LoadDistributions(basedir string) (*Distributions, error) {
	ret := Distributions{}

	file, err := os.Open(basedir + "/conf/distributions")
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)

	for {
		para, err := control.ParseParagraph(reader)
		if err != nil {
			return nil, err
		}
		if para == nil {
			break
		}

		arches := []string{}
		source := false

		for _, arch := range strings.Split(para.Values["Architectures"], " ") {
			if arch == "source" {
				source = true
				continue
			}
			arches = append(arches, arch)
		}

		ret.Blocks = append(ret.Blocks, Distribution{
			Codename:      para.Values["Codename"],
			Suite:         para.Values["Suite"],
			Components:    strings.Split(para.Values["Components"], " "),
			Architectures: arches,
			Source:        source,
			Tracking:      strings.Split(para.Values["Tracking"], " "),
		})
	}

	return &ret, nil
}
