package values

var (
	PathFull                      CloseCode = 3000
	ProtocolError                 CloseCode = 3001
	InternalError                 CloseCode = 3002
	HandoverOfTheSignalingChannel           = 3003
	DroppedByInitiator            CloseCode = 3004
	InitiatorCouldNotDecrypt      CloseCode = 3005
	NoSharedTaskCloseCode                   = 3006
	InvalidKey                    CloseCode = 3007
	Timeout                       CloseCode = 3008

	DropResponderProtocolError            = 3001
	DropResponderInternalError            = 3002
	DropResponderDroppedByInitiator       = 3004
	DropResponderInitiatorCouldNotDecrypt = 3005
)

type CloseCode uint
