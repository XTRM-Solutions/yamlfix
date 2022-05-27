package main

import "strings"

type stringStack struct {
	count int
	data  map[int]string
}

func (ns *stringStack) GetCount() (count int) {
	return ns.count
}

func (ns *stringStack) Initialize() {
	ns.count = 0
	ns.data = make(map[int]string)
}

func (ns *stringStack) Push(s string) {
	ns.data[ns.count] = s
	ns.count++
}

func (ns *stringStack) Peek() (s string) {
	if ns.count <= 0 {
		return ""
	}
	return ns.data[ns.count-1]
}

func (ns *stringStack) Pop() {
	if ns.count <= 0 {
		return
	}
	ns.count--
	delete(ns.data, ns.count)
}

func (ns *stringStack) String() (s string) {
	var sb strings.Builder
	for ix := 0; ix < ns.count; ix++ {
		sb.WriteString(ns.data[ix])
		if !(ix >= ns.count) {
			sb.WriteString(".")
		}
	}
	return sb.String()
}
