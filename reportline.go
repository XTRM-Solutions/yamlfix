package main

import (
	"strings"
)

const CSVSEPCHAR = "\t"

const NEWLINECHAR = "\n"

const APIHEADERS = "OperationID" + CSVSEPCHAR +
	"MediaType" + CSVSEPCHAR +
	"ParameterType" + CSVSEPCHAR +
	"ParameterName" + CSVSEPCHAR +
	"FullyQualifiedParameter" + NEWLINECHAR

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
	ns.count += 1
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
	ns.count -= 1
	delete(ns.data, ns.count)
}

func (ns *stringStack) Concat() (s string) {
	var sb strings.Builder
	for ix := 0; ix < ns.count; ix++ {
		sb.WriteString(ns.data[ix])
		if ix < (ns.count - 1) {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

type YamlReportLine struct {
	OperationID string
	ParamNames  stringStack
	TypeNames   stringStack
	MediaNames  stringStack
}

func (yl *YamlReportLine) GetHeaders() (s string) {
	return "OperationID" + "\t" +
		"MediaType" + "\t" +
		"ParameterName" + "\t" +
		"FullyQualifiedParameter" + "\n"
}

func (yl *YamlReportLine) String() (s string) {
	var sb strings.Builder
	WriteSB(&sb, yl.OperationID, CSVSEPCHAR,
		yl.MediaNames.Peek(), CSVSEPCHAR,
		yl.TypeNames.Peek(), CSVSEPCHAR,
		yl.ParamNames.Peek(), CSVSEPCHAR,
		yl.ParamNames.Concat(), "\n")
	return sb.String()
}

func (yl *YamlReportLine) Reset() {
	yl.OperationID = ""
	yl.ParamNames.Initialize()
	yl.TypeNames.Initialize()
	yl.MediaNames.Initialize()
}
