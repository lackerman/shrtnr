package utils

import (
	"fmt"
	"hash/fnv"
)

// EncodeURL takes a Url and returns a shortened version of it
func EncodeURL(s string) (string, error) {
	h := fnv.New32a()
	if _, err := h.Write([]byte(s)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum32()), nil
}
