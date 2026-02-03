package traffic

import "time"

type CaptureID string

type CapturedTraffic struct {
	ID        CaptureID
	CapturedAt time.Time

	Method string
	URL    string
}
