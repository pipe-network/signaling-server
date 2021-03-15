package values

import (
	"bytes"
	"encoding/hex"
	"errors"
)

var (
	HexKeyNot32BytesLong = errors.New("given hexKey is not 32 bytes (64 hex chars) long")
)

type Key [32]byte

func FromHex(hexKey string) (*Key, error) {
	decodedKey, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}

	if len(decodedKey) != len(Key{}) {
		return nil, HexKeyNot32BytesLong
	}

	key := Key{}
	copy(key[:], decodedKey[:32])

	return &key, err

}

func (k Key) Bytes() [32]byte {
	return k
}

func (k Key) Empty() bool {
	emptyKey := Key{}
	return bytes.Equal(k[:], emptyKey[:])
}

func (k Key) Equals(key Key) bool {
	return bytes.Equal(k[:], key[:])
}
