package trafficlight

type lightStatus int

const (
	LIGHT_STATUS_RED lightStatus = iota
	LIGHT_STATUS_GREEN
	LIGHT_STATUS_YELLOW
)

// Helper to print light status (optional)
func (s lightStatus) String() string {
	switch s {
	case LIGHT_STATUS_RED:
		return "RED"
	case LIGHT_STATUS_GREEN:
		return "GREEN"
	case LIGHT_STATUS_YELLOW:
		return "YELLOW"
	default:
		return "UNKNOWN"
	}
}
