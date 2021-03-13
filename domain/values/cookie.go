package values

import "crypto/rand"

type Cookie [16]byte

func NewRandomCookie() (*Cookie, error) {
	cookie := Cookie{}
	_, err := rand.Read(cookie[:])
	if err != nil {
		return nil, err
	}
	return &cookie, nil
}
