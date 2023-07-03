// eventmanager_singleton.go
package event

var eventManagerInstance *EventManager

// GetEventManagerInstance returns the singleton instance of EventManager.
func GetEventManagerInstance() *EventManager {
	if eventManagerInstance == nil {
		eventManagerInstance = NewEventManager()
	}
	return eventManagerInstance
}
