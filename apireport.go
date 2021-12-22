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
	// LIFO order for defer
	defer DeferError(outFile.Close)
	defer DeferError(outWriter.Flush)
	_, err = outWriter.WriteString(APIHEADERS)
	if nil != err {
		xLog.Fatalf(
			"Failed to write API Headers to %s because %s\n",
			apiFileName, err.Error())
	}

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

	yl.Reset()

	if nil == item {
		return
	}

	yl.OperationID = item.OperationID
	if FlagDebug {
		xLog.Printf("operation id: %s\n", item.OperationID)
	}
	doContent(&yl, item.RequestBody.Value.Content)
}

func doContent(yl *YamlReportLine, c oas.Content) {
	if nil == yl || nil == c {
		return
	}

	for key := range c {
		yl.MediaNames.Push(key)
		ref, ok := c[key]
		if ok && nil != ref && nil != ref.Schema {
			doSchema(yl, ref.Schema.Value)
		}
		yl.MediaNames.Pop()
	}

}

func doSchemas(yl *YamlReportLine, schemas oas.Schemas) {
	if nil == yl || nil == schemas {
		return
	}
	for key := range schemas {
		ref, ok := schemas[key]
		if ok && nil != ref && nil != ref.Value {
			yl.ParamNames.Push(key)
			doSchema(yl, ref.Value)
			yl.ParamNames.Pop()
		}
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

	yl.TypeNames.Push(schema.Type)
	_, err := outWriter.WriteString(yl.String())
	if err != nil {
		xLog.Fatalf("outWriter.WriteString(\"%s\") failed because %s\n",
			yl.String(), err.Error())
	}
	if FlagDebug {
		xLog.Println(yl.String())
	}

	doSchemaRefs(yl, &schema.OneOf)
	doSchemaRefs(yl, &schema.AnyOf)
	doSchemaRefs(yl, &schema.AllOf)
	if nil != schema.Not {
		doSchema(yl, schema.Not.Value)
	}
	if nil != schema.Properties {
		doSchemas(yl, schema.Properties)
	}
	if nil != schema.Items {
		doSchema(yl, schema.Items.Value)
	}
	yl.TypeNames.Pop()
}
