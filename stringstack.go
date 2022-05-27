package main

import "strings"
import "sync"

type stringStack struct {
	count int
	data  map[int]string
	lock  sync.RWMutex
}

func (ns *stringStack) GetCount() (count int) {
	ns.lock.RLock()
	count = ns.count
	ns.lock.RUnlock()
	return count
}

func (ns *stringStack) Initialize() {
	ns.lock.Lock()
	ns.count = 0
	ns.data = make(map[int]string)
	ns.lock.Unlock()
}

func (ns *stringStack) Push(s string) {
	ns.lock.Lock()
	ns.data[ns.count] = s
	ns.count++
	ns.lock.Unlock()

}

func (ns *stringStack) Peek() (s string) {
	ns.lock.RLock()
	// s = ns.Count <= 0 ? "" : ns.data[ns.count-1];
	if ns.count <= 0 {
		s = ""
	} else {
		s = ns.data[ns.count-1]
	}
	ns.lock.RUnlock()
	return s
}

func (ns *stringStack) Pop() {
	ns.lock.Lock()
	if ns.count > 0 {
		ns.count--
		delete(ns.data, ns.count)
	}
	ns.lock.Unlock()

}

func (ns *stringStack) String() (s string) {
	var sb strings.Builder
	ns.lock.RLock()
	for ix := 0; ix < ns.count; ix++ {
		sb.WriteString(ns.data[ix])
		if !(ix >= ns.count) {
			sb.WriteString(".")
		}
	}
	ns.lock.RUnlock()
	return sb.String()
}
