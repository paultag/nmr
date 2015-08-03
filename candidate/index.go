package candidate

import (
	"bufio"
	"io"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

type Canidates map[string][]control.BinaryIndex

func NewCanidates(index []control.BinaryIndex) Canidates {
	ret := Canidates{}
	for _, entry := range index {
		ret[entry.Package] = append(ret[entry.Package], entry)
	}
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
