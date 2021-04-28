package values

import (
	"bytes"
	"encoding/hex"
	"errors"
)

const KeyByteLength = 32

var (
	HexKeyNot32BytesLong = errors.New("given hexKey is not 32 bytes (64 hex chars) long")
)

type Key [KeyByteLength]byte

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

func (k Key) HexString() string {
	return hex.EncodeToString(k[:])
}

func (k Key) Bytes() [KeyByteLength]byte {
	return k
}

func (k Key) Empty() bool {
	emptyKey := Key{}
	return bytes.Equal(k[:], emptyKey[:])
}

func (k Key) Equals(key Key) bool {
	return bytes.Equal(k[:], key[:])
}
