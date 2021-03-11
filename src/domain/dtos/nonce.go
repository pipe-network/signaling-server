package dtos

type Nonce struct {
	Cookie         Cookie
	Source         Address
	Destination    Address
	CombinedSequenceNumber CombinedSequenceNumber
}
