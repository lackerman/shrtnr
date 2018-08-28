package utils

import (
	"fmt"
	"hash/fnv"
)

// EncodeURL takes a Url and returns a shortened version of it
func EncodeURL(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum32())
}
