package main

import (
	"bufio"
	oas "github.com/getkin/kin-openapi/openapi3"
	"os"
)

var outWriter *bufio.Writer

func ApiReport(api *oas.T) {
	const apiFileName = "apireport.csv"

	if nil == api {
		return
	}

	outFile, err := os.Create(apiFileName)
	if nil != err {
		xLog.Fatalf("Attempt to open file %s failed because %s",
			apiFileName, err.Error())
	}
	outWriter = bufio.NewWriter(outFile)
	DeferError(outWriter.Flush)
	DeferError(outFile.Close)

	for _, val01 := range api.Paths {
		apiCallReport(val01)
	}
}

func apiCallReport(item *oas.PathItem) {
	operationParamReport(item.Connect)
	operationParamReport(item.Delete)
	operationParamReport(item.Get)
	operationParamReport(item.Head)
	operationParamReport(item.Options)
	operationParamReport(item.Patch)
	operationParamReport(item.Post)
	operationParamReport(item.Put)
	operationParamReport(item.Trace)
}

func operationParamReport(item *oas.Operation) {
	var yl YamlReportLine

	if nil == item {
		return
	}

	yl.opID = item.OperationID
	if FlagDebug {
		xLog.Printf("operation id: %s\n", item.OperationID)
	}
	doContent(&yl, item.RequestBody.Value.Content)

}

func doContent(yl *YamlReportLine, content oas.Content) {
	if nil == content {
		return
	}
	for _, c := range content {
		doSchema(yl, c.Schema.Value)
	}
}

func doSchemaRefs(yl *YamlReportLine, sr *oas.SchemaRefs) {
	if nil == yl || nil == sr {
		return
	}
	for _, schemaRef := range *sr {
		doSchema(yl, schemaRef.Value)
	}

}

func doSchema(yl *YamlReportLine, schema *oas.Schema) {
	if nil == schema || nil == yl {
		return
	}
	yl.PushType(*&schema.Type)
	doSchemaRefs(yl, &schema.OneOf)
	doSchemaRefs(yl, &schema.AnyOf)
	doSchemaRefs(yl, &schema.AllOf)
	if nil != schema.Not {
		doSchema(yl, schema.Not.Value)
	}
	if nil != schema.Properties {
		for _, s := range schema.Properties {
			doSchema(yl, s.Value)
		}
	}
	_ = yl.PopType()
}
