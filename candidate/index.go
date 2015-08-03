package candidate

import (
	"bufio"
	"fmt"
	"io"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

type Canidates map[string][]control.BinaryIndex

func (can *Canidates) AppendBinaryIndexReader(in io.Reader) error {
	reader := bufio.NewReader(in)
	index, err := control.ParseBinaryIndex(reader)
	if err != nil {
		return err
	}
	can.AppendBinaryIndex(index)
	return nil
}

func (can *Canidates) AppendBinaryIndex(index []control.BinaryIndex) {
	for _, entry := range index {
		(*can)[entry.Package] = append((*can)[entry.Package], entry)
	}
}

func NewCanidates(index []control.BinaryIndex) Canidates {
	ret := Canidates{}
	ret.AppendBinaryIndex(index)
	return ret
}

func ReadFromBinaryIndex(in io.Reader) (*Canidates, error) {
	reader := bufio.NewReader(in)
	index, err := control.ParseBinaryIndex(reader)
	if err != nil {
		return nil, err
	}
	can := NewCanidates(index)
	return &can, nil
}

func (can Canidates) ExplainSatisfiesBuildDepends(arch dependency.Arch, depends dependency.Dependency) (bool, string) {
	for _, possi := range depends.GetPossibilities(arch) {
		if !can.Satisfies(possi) {
			return false, fmt.Sprintf("Possi %s can't be satisfied.", possi.Name)
		}
	}
	return true, "All relations are a go"
}

func (can Canidates) SatisfiesBuildDepends(arch dependency.Arch, depends dependency.Dependency) bool {
	ret, _ := can.ExplainSatisfiesBuildDepends(arch, depends)
	return ret
}

func (can Canidates) Satisfies(possi dependency.Possibility) bool {
	///
	///  XXX: DON'T IGNORE ARCHES
	///

	entries, ok := can[possi.Name]
	if !ok { // no known entries in the Index
		return false
	}

	if possi.Version == nil {
		return true
	}

	// OK, so we have to play with versions now.
	vr := *possi.Version
	relatioNumber, _ := version.Parse(vr.Number)

	for _, installable := range entries {
		q := version.Compare(installable.Version, relatioNumber)

		switch vr.Operator {
		case ">=":
			return q >= 0
		case "<=":
			return q <= 0
		case ">>":
			return q > 0
		case "<<":
			return q < 0
		case "=":
			return q == 0
		default:
			return false // XXX: WHAT THE SHIT
		}
	}

	return false
}
