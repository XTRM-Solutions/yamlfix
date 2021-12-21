package main

import "strings"

const SEPARATORCHAR = "\t"

type StringStack struct {
	count   int
	strings []string
}

func (ns *StringStack) GetCount() (count int) {
	return ns.count
}

func (ns *StringStack) Reset() {
	ns.count = 0
	ns.strings = nil
}

func (ns *StringStack) Push(s string) {
	ns.strings[ns.count] = s
	ns.count += 1
}

func (ns *StringStack) Peek() (s string) {
	if ns.count <= 0 {
		return ""
	}
	return ns.strings[ns.count-1]
}

func (ns *StringStack) Pop() {
	if ns.count <= 0 {
		return
	}
	ns.count -= 1
}

func (ns *StringStack) Concat() (s string) {
	var sb strings.Builder
	for ix := 0; ix < ns.count; ix++ {
		sb.WriteString(ns.strings[ix])
		if ix < (ns.count - 1) {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

type YamlReportLine struct {
	OperationID string
	ParamNames  StringStack
	TypeNames   StringStack
	MediaNames  StringStack
}

func (yl *YamlReportLine) String() (s string) {
	var sb strings.Builder
	WriteSB(&sb, yl.OperationID, SEPARATORCHAR,
		yl.MediaNames.Peek(), SEPARATORCHAR,
		yl.TypeNames.Peek(), SEPARATORCHAR,
		yl.ParamNames.Concat())
	return sb.String()
}

func (yl *YamlReportLine) Reset() {
	yl.OperationID = ""
	yl.ParamNames.Reset()
	yl.TypeNames.Reset()
	yl.MediaNames.Reset()
}
