package constants

const (
	Registered = "REGISTERED"
	Processing = "PROCESSING"
	Invalid    = "INVALID"
	Processed  = "PROCESSED"
)

func GetStatusOrder(status string) string {
	switch status {
	case Registered:
		return "NEW"
	case Processing:
		return "PROCESSING"
	case Invalid:
		return "INVALID"
	case "PROCESSED":
		return "PROCESSED"
	default:
		return "UNSUPPORTED"
	}
}
