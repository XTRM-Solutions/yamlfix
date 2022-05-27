package main

import "strings"

type YamlReportLine struct {
	OperationID string
	ParamNames  stringStack
	TypeNames   stringStack
	MediaNames  stringStack
}

func (yl *YamlReportLine) String() (s string) {
	var sb strings.Builder
	WriteSB(&sb, yl.OperationID, CSVSEPCHAR,
		yl.MediaNames.Peek(), CSVSEPCHAR,
		yl.TypeNames.Peek(), CSVSEPCHAR,
		yl.ParamNames.Peek(), CSVSEPCHAR,
		yl.ParamNames.String(), "\n")
	return sb.String()
}

func (yl *YamlReportLine) Reset() {
	yl.OperationID = ""
	yl.ParamNames.Initialize()
	yl.TypeNames.Initialize()
	yl.MediaNames.Initialize()
}
