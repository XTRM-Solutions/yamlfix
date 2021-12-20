package main

import (
	"bufio"
	"github.com/NanXiao/stack"
	oas "github.com/getkin/kin-openapi/openapi3"
	"os"
	"strings"
	// https://pkg.go.dev/github.com/getkin/kin-openapi@v0.85.0/openapi3
)

var outWriter *bufio.Writer
var paramStack stack.Stack
var ps = &paramStack

func ApiReport(api *oas.T) {
	const apiFileName = "apireport.csv"

	if nil == api {
		return
	}

	outFile, err := os.Create(apiFileName)
	outWriter = bufio.NewWriter(outFile)
	if nil != err {
		xLog.Fatalf("Attempt to open file %s failed because",
			apiFileName, err.Error())
	}
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
	if nil == item {
		return
	}

	if IsStringSet(&item.OperationID) {
		ps.Push(item.OperationID)
	} else {
		ps.Push("Unspecified OperationID")
	}
	xLog.Printf("operation id: %s\n", item.OperationID)

	var r strings.Builder
	WriteSB(&r, "REQUEST: %s\t", item.OperationID)
	doContent(&r, item.RequestBody.Value.Content)

	_ = ps.Pop()
}

func doContent(r *strings.Builder, content oas.Content) {
	if nil == content {
		return
	}
	for _, c := range content {
		doSchema(r, c.Schema.Value)
	}
}

func doSchema(r *strings.Builder, schema *oas.Schema) {
	if nil == schema || nil == r {
		return
	}

}
