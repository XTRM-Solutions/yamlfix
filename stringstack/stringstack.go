package stringstack

import (
	"strings"
	"sync"
)

type StringStack struct {
	count int
	data  map[int]string
	lock  sync.RWMutex
}

func (ns *StringStack) Rlock() {
	ns.lock.RLock()
}
func (ns *StringStack) RUnlock() {
	ns.lock.RUnlock()
}
func (ns *StringStack) Lock() {
	ns.lock.Lock()
}
func (ns *StringStack) Unlock() {
	ns.lock.Unlock()
}

func (ns *StringStack) GetCount() (count int) {
	ns.Rlock()
	{
		count = ns.count
	}
	ns.RUnlock()
	return count
}

func (ns *StringStack) Initialize() {
	ns.Lock()
	{
		ns.count = 0
		ns.data = make(map[int]string)
	}
	ns.Unlock()
}

func (ns *StringStack) Push(s string) {
	ns.Lock()
	{
		ns.data[ns.count] = s
		ns.count++
	}
	ns.RUnlock()

}

func (ns *StringStack) Peek() (s string) {
	ns.lock.RLock()
	{
		// s = ns.Count <= 0 ? "" : ns.data[ns.count-1];
		if ns.count <= 0 {
			s = ""
		} else {
			s = ns.data[ns.count-1]
		}
		ns.lock.RUnlock()
	}
	return s
}

func (ns *StringStack) Pop() {
	ns.lock.Lock()
	{
		if ns.count > 0 {
			ns.count--
			delete(ns.data, ns.count)
		}
	}
	ns.lock.Unlock()
}

func (ns *StringStack) String() (s string) {
	var sb strings.Builder
	ns.lock.RLock()
	{
		for ix := 0; ix < ns.count; ix++ {
			sb.WriteString(ns.data[ix])
			if !(ix >= ns.count) {
				sb.WriteString(".")
			}
		}
	}
	ns.lock.RUnlock()
	return sb.String()
}
