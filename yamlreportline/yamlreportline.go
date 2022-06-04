package yamlreportline

import (
	"strings"
	"yamlfix/misc"
	"yamlfix/stringstack"
)

var csvSepChar string = "\t"

type YamlReportLine struct {
	OperationID string
	ParamNames  stringstack.StringStack
	TypeNames   stringstack.StringStack
	MediaNames  stringstack.StringStack
}

func (yl *YamlReportLine) String() (s string) {
	var sb strings.Builder
	misc.WriteSB(&sb, yl.OperationID, csvSepChar,
		yl.MediaNames.Peek(), csvSepChar,
		yl.TypeNames.Peek(), csvSepChar,
		yl.ParamNames.Peek(), csvSepChar,
		yl.ParamNames.String(), "\n")
	return sb.String()
}

func (yl *YamlReportLine) Reset() {
	yl.OperationID = ""
	yl.ParamNames.Initialize()
	yl.TypeNames.Initialize()
	yl.MediaNames.Initialize()
}

func (yl *YamlReportLine) SetSeparatorChar(sepchar string) {
	csvSepChar = sepchar
}
