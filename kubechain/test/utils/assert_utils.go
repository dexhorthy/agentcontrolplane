package utils

import (
	"strings"
	"time"

	. "github.com/onsi/gomega"

	"k8s.io/client-go/tools/record"
)

type eventAssertion struct {
	eventRecorder *record.FakeRecorder
}

// ExpectEvent starts the fluent assertion chain for event recording
func ExpectEvent(recorder *record.FakeRecorder) *eventAssertion {
	return &eventAssertion{
		eventRecorder: recorder,
	}
}

// ToEmitEventContaining completes the fluent assertion chain and checks for the event
func (a *eventAssertion) ToEmitEventContaining(substring string) {
	Eventually(func() bool {
		select {
		case event := <-a.eventRecorder.Events:
			return strings.Contains(event, substring)
		default:
			return false
		}
	}, 5*time.Second, 100*time.Millisecond).Should(BeTrue())
}
