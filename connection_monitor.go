package actioncable
/* TODO: this unsed code */

import (
	"time"
	"log"
)

const (
	// use the same value as actioncable.js
	// 1 ping is 2 sec interval, so detect stale when 2 ping missing.
	//DEFAULT_STALE_THRESHOLD time.Duration = 6 * time.Second
)

type connectionMonitor struct {
	connection *connection
	reconnectAttempts int
	pingedAt  time.Time
	startedAt time.Time
	stopedAt  time.Time
	disconnectedAt time.Time

	isPolling bool
	stopPollingCh chan struct{}

	staleThreshold time.Duration
}

func newConnectionMonitor() *connectionMonitor {
	return &connectionMonitor{
		isPolling: false,
		stopPollingCh: make(chan struct{}),
		staleThreshold: DEFAULT_STALE_THRESHOLD,
	}
}

func (m *connectionMonitor) start() {
	if (m.isRunning()) {
		return
	}

	m.startedAt = time.Now()
	m.stopedAt  = time.Time{}
	go m.startPolling()
}

func (m *connectionMonitor) stop() {
}

func (m *connectionMonitor) isRunning() bool {
	return ! m.startedAt.IsZero() && m.stopedAt.IsZero()
}

func (m *connectionMonitor) recordConnect() {
	m.recordPing()
	m.reconnectAttempts = 0
	m.disconnectedAt = time.Time{}
}

func (m *connectionMonitor) recordPing() {
	m.pingedAt = time.Now()
}

func (m *connectionMonitor) startPolling() {
	m.stopPolling()
	m.poll()
}

func (m *connectionMonitor) poll() {
	defer (func() { m.isPolling = false })()

	for {
		select {
		case <-time.After(m.staleThreshold):
			if m.connectionIsStale() {
				log.Println("connection is stale")
				// call stop
				//reconnectCh <-
				return
			}
			log.Println("not stale")
		case _ = <-m.stopPollingCh:
			// logging
			return
		}
	}
}

func (m *connectionMonitor) connectionIsStale() bool {
	threshold := time.Now().Add(-m.staleThreshold)

	log.Println("ok")
	log.Println(threshold)
	log.Println(m.pingedAt)

	return threshold.After(m.pingedAt)
}

func (m *connectionMonitor) stopPolling() {
	if m.isPolling {
		m.stopPollingCh <-struct{}{}
	}
}
