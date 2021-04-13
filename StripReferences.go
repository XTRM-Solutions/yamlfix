package main

import (
	oas "github.com/getkin/kin-openapi/openapi3"
	// https://pkg.go.dev/github.com/getkin/kin-openapi@v0.53.0/openapi3
)

// StripReferences Remove the reference pointers so that
// UnMarshal expands everything.

// StripReferences
// Remove the reference pointer from the loaded OAS data
// thus, when the OAS data is unmarshalled back into JSON
// there are no internal references.
func StripReferences(api *oas.Swagger) {
	if nil == api {
		return
	}
	for _, val01 := range api.Paths {
		val01.Ref = ""
		StripPathItem(val01)
	}
}

// StripOperationReferences
// For a given oas.Operation, remove the references
// within it (not really a public function)
func StripOperationReferences(op *oas.Operation) {
	if nil == op {
		return
	}
	op.RequestBody.Ref = ""
	for _, val02 := range op.RequestBody.Value.Content {
		StripReferencesSchema(val02.Schema)
	}
	for _, val03 := range op.Responses {
		val03.Ref = ""
		for _, val04 := range val03.Value.Content {
			StripReferencesSchema(val04.Schema)
		}
	}
	for _, val05 := range op.Callbacks {
		val05.Ref = ""
		for _, val06 := range *val05.Value {
			StripPathItem(val06)
		}
	}
}

// StripPathItem
// For a given oas.PathItem (an Operation), remove the references
// within it (not really a public function)
func StripPathItem(item *oas.PathItem) {
	if nil == item {
		return
	}
	StripOperationReferences(item.Connect)
	StripOperationReferences(item.Delete)
	StripOperationReferences(item.Get)
	StripOperationReferences(item.Head)
	StripOperationReferences(item.Options)
	StripOperationReferences(item.Patch)
	StripOperationReferences(item.Post)
	StripOperationReferences(item.Put)
	StripOperationReferences(item.Trace)
}

// StripReferencesSchema Eventually, whatever reference
// paths an OAS file has comes down to schema references,
// and this recursively clears those references
func StripReferencesSchema(schema *oas.SchemaRef) {
	if nil == schema {
		return
	}
	// clean this reference, and look for sub-references
	// within the schemaBody
	schema.Ref = ""
	schemaBody := schema.Value
	if 0 != len(schemaBody.Extensions) {
		xLog.Print("WARNING: Extension Properties references are NOT handled by this program")
	}
	if nil != schemaBody.OneOf {
		for _, val01 := range schemaBody.OneOf {
			StripReferencesSchema(val01)
		}
	}
	if nil != schemaBody.AnyOf {
		for _, val02 := range schemaBody.AnyOf {
			StripReferencesSchema(val02)
		}
	}
	if nil != schemaBody.AllOf {
		for _, val03 := range schemaBody.AllOf {
			StripReferencesSchema(val03)
		}
	}
	if nil != schemaBody.Not {
		StripReferencesSchema(schemaBody.Not)
	}
	if nil != schemaBody.Items {
		StripReferencesSchema(schema.Value.Items)
	}
	if nil != schemaBody.Properties {
		for _, val04 := range schemaBody.Properties {
			StripReferencesSchema(val04)
		}
	}
	if nil != schemaBody.AdditionalProperties {
		StripReferencesSchema(schemaBody.AdditionalProperties)
	}
}
