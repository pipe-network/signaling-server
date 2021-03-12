package domain_services

import "math/big"

func BytesToInteger(bytes []byte) int {
	number := big.Int{}
	number.SetBytes(bytes[:])
	return int(number.Int64())
}
