// eventmanager.go
package event

import (
	"sort"
	"time"
)

type EventManager struct {
	database []EventRecord
}

func NewEventManager() *EventManager {
	return &EventManager{}
}

func (em *EventManager) AddEventToDatabase(event EventRecord) {
	em.database = append(em.database, event)
}

func (em *EventManager) SortEvents(userID int64) []EventRecord {
	var futureEvents []EventRecord

	currentTime := time.Now()

	for _, event := range em.database {
		if event.UserID == userID && event.EventTime.After(currentTime) {
			futureEvents = append(futureEvents, event)
		}
	}

	sort.Slice(futureEvents, func(i, j int) bool {
		return futureEvents[i].EventTime.Before(futureEvents[j].EventTime)
	})

	return futureEvents
}
