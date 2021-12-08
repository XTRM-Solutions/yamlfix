package main

/*************************
This all works very simply because the openapi3 toolkit tracks
(and creates and removes) references in a 2-layer ref/value
structure. A NIL ref means that this is NOT a reference; a
non-NIL reference means, create this as the specified reference.
Therefore, to 'strip' references, rather than copy & pasting
trees and maps and submaps into maps and trees and yet more
submaps, one need only clear the Ref string.
SO MUCH EASY! THANK YOU OPENAPI LIBRARY!
*/

import (
	oas "github.com/getkin/kin-openapi/openapi3"
	// https://pkg.go.dev/github.com/getkin/kin-openapi@v0.53.0/openapi3
)

// StripReferences
// Remove the reference pointer from the loaded OAS data
// thus, when the OAS data is unmarshalled back into JSON
// there are no internal references.
func StripReferences(api *oas.T) {
	if nil == api {
		return
	}
	for _, val01 := range api.Paths {
		val01.Ref = ""
		StripPathItem(val01)
	}
}

// StripPathItem
// For a given oas.PathItem, remove the references
// within it. Since all the references occur within
// an Operation ... just clear the Operations
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

// StripOperationReferences
// For a given oas.Operation, remove the references
// within it. That consists of cleaning the Schema
// references and callbacks (which are Pathitems)
func StripOperationReferences(op *oas.Operation) {
	if nil == op {
		return
	}
	op.RequestBody.Ref = ""
	for _, val01 := range op.RequestBody.Value.Content {
		StripReferencesSchema(val01.Schema)
	}
	for _, val02 := range op.Responses {
		val02.Ref = ""
		for _, val03 := range val02.Value.Content {
			StripReferencesSchema(val03.Schema)
		}
	}
	for _, val04 := range op.Callbacks {
		val04.Ref = ""
		for _, val05 := range *val04.Value {
			StripPathItem(val05)
		}
	}
}

// StripReferencesSchema Eventually, whatever reference
// paths an OAS file has comes down to schema references,
// and this recursively clears those schemas as well as
// all the schema references within
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
