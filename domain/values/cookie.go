package values

import (
	"bytes"
	"crypto/rand"
)

type Cookie [16]byte

func NewRandomCookie() (*Cookie, error) {
	cookie := Cookie{}
	_, err := rand.Read(cookie[:])
	if err != nil {
		return nil, err
	}
	return &cookie, nil
}

func (c Cookie) Equal(cookie Cookie) bool {
	return bytes.Equal(cookie[:], c[:])
}

func (c Cookie) Empty() bool {
	return c.Equal(Cookie{})
}
