// eventrecord.go
package event

import "time"

type EventRecord struct {
	UserID    int64
	EventTime time.Time
	RawData   string
}
