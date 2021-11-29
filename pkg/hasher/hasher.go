package hasher

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
)

type Hasher interface {
	ComputeHash(str string) string
}

type sha1Hasher struct {
	hash hash.Hash
}

func NewSha1Hasher() Hasher {
	return &sha1Hasher{
		sha1.New(),
	}
}

func (h *sha1Hasher) ComputeHash(str string) string {
	b := h.hash.Sum([]byte(str))
	res := hex.EncodeToString(b)
	return res
}
