package helpers

import (
	"encoding/hex"
	"hash"
	"io"
	"os"
)

func HashFile(path string, algo hash.Hash) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(algo, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(algo.Sum(nil)), nil
}
