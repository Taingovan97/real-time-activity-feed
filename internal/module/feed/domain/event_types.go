// Package domain provides domain entities and constants for the activity feed module.
package domain

import (
	"fmt"
	"strings"
)

// Fixed feed event types accepted by the system.
const (
	EventTypeLogin        = "login"
	EventTypeUpload       = "upload"
	EventTypeComment      = "comment"
	EventTypeNotification = "notification"
	EventTypeReport       = "report"
	EventTypeTask         = "task"
	EventTypeMessage      = "message"
	EventTypeSync         = "sync"
	EventTypeStatus       = "status"
	EventTypeApproval     = "approval"
)

// AllowedEventTypes is the fixed list of event types the system accepts.
var AllowedEventTypes = []string{
	EventTypeLogin,
	EventTypeUpload,
	EventTypeComment,
	EventTypeNotification,
	EventTypeReport,
	EventTypeTask,
	EventTypeMessage,
	EventTypeSync,
	EventTypeStatus,
	EventTypeApproval,
}

// InvalidEventTypeError is returned when callers use an unsupported event type.
type InvalidEventTypeError struct {
	Value string
}

// Error implements the error interface.
func (e *InvalidEventTypeError) Error() string {
	return fmt.Sprintf("event_type must be one of: %s", strings.Join(AllowedEventTypes, ", "))
}

// IsValidEventType reports whether the provided event type is supported.
func IsValidEventType(eventType string) bool {
	for _, allowed := range AllowedEventTypes {
		if eventType == allowed {
			return true
		}
	}
	return false
}

// ValidateEventType returns an error if the provided event type is unsupported.
func ValidateEventType(eventType string) error {
	if eventType == "" || IsValidEventType(eventType) {
		return nil
	}

	return &InvalidEventTypeError{Value: eventType}
}
