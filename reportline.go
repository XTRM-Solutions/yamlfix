package main

import (
	"github.com/NanXiao/stack"
	"strings"
)

type YamlReportLine struct {
	opID       string
	paramCount int
	paramName  []string
	paramType  stack.Stack
}

func (yl *YamlReportLine) PushType(s string) {
	yl.paramType.Push(s)
}
func (yl *YamlReportLine) PopType() (s string) {
	return yl.paramType.Pop().(string)
}
func (yl *YamlReportLine) PeekType() (s string) {
	return yl.paramType.Top().(string)
}

func (yl *YamlReportLine) PushParam(s string) {
	yl.paramName[yl.paramCount] = s
	yl.paramCount += 1
}

func (yl *YamlReportLine) PopParam() {
	if 0 >= yl.paramCount {
		return
	}
	yl.paramCount -= 1
}

func (yl *YamlReportLine) String() (s string) {
	var rLine strings.Builder
	WriteSB(&rLine, yl.opID, "\t", yl.PeekType(), "\t")

	if nil != yl.paramName {
		for ix := 0; ix < yl.paramCount; ix++ {
			rLine.WriteString(yl.paramName[ix])
			if ix < yl.paramCount-1 {
				rLine.WriteString(".")
			}
		}

	}
	return rLine.String()
}
