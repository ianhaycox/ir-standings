// Package events fake windows events to ease Linux development
package events

import "time"

func OpenEvent(eventName string) {
}

func WaitForSingleObject(timeout time.Duration) bool {
	return false
}

func BroadcastMsg(msgName string, msg int, p1 int, p2 interface{}, p3 int) bool {
	return false
}
