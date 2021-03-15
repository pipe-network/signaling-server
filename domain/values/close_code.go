package values

var (
	PathFullCode                      CloseCode = 3000
	ProtocolErrorCode                 CloseCode = 3001
	InternalErrorCode                 CloseCode = 3002
	HandoverOfTheSignalingChannelCode CloseCode = 3003
	DroppedByInitiatorCode            CloseCode = 3004
	InitiatorCouldNotDecryptCode      CloseCode = 3005
	NoSharedTaskFoundCode             CloseCode = 3006
	InvalidKeyCode                    CloseCode = 3007
	TimeoutCode                       CloseCode = 3008

	DropResponderProtocolErrorCode            CloseCode = 3001
	DropResponderInternalErrorCode            CloseCode = 3002
	DropResponderDroppedByInitiatorCode       CloseCode = 3004
	DropResponderInitiatorCouldNotDecryptCode CloseCode = 3005
)

type CloseCode uint

func (c CloseCode) Message() string {
	switch c {
	case PathFullCode:
		return "Path Full"
	case ProtocolErrorCode:
		return "Protocol Error"
	case InternalErrorCode:
		return "Internal Error"
	case HandoverOfTheSignalingChannelCode:
		return "Handover of the Signalling Channel"
	case DroppedByInitiatorCode:
		return "Dropped by Initiator"
	case InitiatorCouldNotDecryptCode:
		return "Initiator Could Not Decrypt"
	case NoSharedTaskFoundCode:
		return "No Shared Task Found"
	case InvalidKeyCode:
		return "Invalid Key"
	case TimeoutCode:
		return "Timeout"
	case DropResponderProtocolErrorCode:
		return "Protocol Error"
	case DropResponderInternalErrorCode:
		return "Internal Error"
	case DropResponderDroppedByInitiatorCode:
		return "Dropped by Initiator"
	case DropResponderInitiatorCouldNotDecryptCode:
		return "Initiator Could Not Decrypt"
	}
	return ""
}

func (c CloseCode) Int() int {
	return int(c)
}
