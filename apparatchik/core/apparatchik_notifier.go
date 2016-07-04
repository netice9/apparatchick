/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/gogen/specific
* THIS FILE SHOULD NOT BE EDITED BY HAND
*/

package core

import "sync"

// ApparatchikNotifier is a event broadcaster
type ApparatchikNotifier struct {
	sync.Mutex
	listeners		[]chan []string
	lastNotification	[]string
}

// NewApparatchikNotifier creates a new ApparatchikNotifier with initial notfication value
func NewApparatchikNotifier(firstNotification []string) *ApparatchikNotifier {
	return &ApparatchikNotifier{
		lastNotification: firstNotification,
	}
}

// Notify notifies current value to all listeners
func (n *ApparatchikNotifier) Notify(value []string) {

	nonBlockingSendToChannel := func(chn chan []string, val []string) {
		// recover in the case of sending to closed channel
		defer func() {
			if r := recover(); r != nil {
				// ignore?
			}
		}()

		select {
		case chn <- val:
			// everything is ok
		default:
			// previous value is blocking the channel, remove it!
			select {
			case <-chn:
				// removed value, all clear to send!
				chn <- val
			default:
				// receiver read it, send it now!
				chn <- val
			}
		}

	}

	n.Lock()
	defer n.Unlock()

	n.lastNotification = value
	for _, listener := range n.listeners {
		nonBlockingSendToChannel(listener, value)
	}

}

// AddListener creats a new listener channel
func (n *ApparatchikNotifier) AddListener(capacity int) chan []string {
	if capacity == 0 {
		capacity = 1
	}
	listenerChannel := make(chan []string, capacity)
	n.Lock()
	defer n.Unlock()
	n.listeners = append(n.listeners, listenerChannel)
	last := n.lastNotification
	listenerChannel <- last
	return listenerChannel
}

// RemoveListener removes and closes an existing listener channel
func (n *ApparatchikNotifier) RemoveListener(listenerChannel chan []string) {
	n.Lock()
	defer n.Unlock()
	filtered := []chan []string{}
	for _, existing := range n.listeners {
		if existing != listenerChannel {
			filtered = append(filtered, existing)
		}
	}
	n.listeners = filtered
	close(listenerChannel)
}

// Close closes and removes all listeners
func (n *ApparatchikNotifier) Close() {
	n.Lock()
	defer n.Unlock()
	for _, listener := range n.listeners {
		close(listener)
	}
	n.listeners = []chan []string{}
}

// NumberOfListeners returns the current count of open listeners
func (n *ApparatchikNotifier) NumberOfListeners() int {
	n.Lock()
	defer n.Unlock()
	return len(n.listeners)
}
