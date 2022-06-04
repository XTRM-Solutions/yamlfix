package main

import (
	oas "github.com/getkin/kin-openapi/openapi3"
	"runtime"
	"strings"
	"yamlfix/misc"
)

// SimplexEnhanceDescriptions
// Find the lower-level parameters, and build a
// table of descriptions
func SimplexEnhanceDescriptions(api *oas.T) {

	// this strips the required column from the response table
	// without requiring 2 different getSchemaProperties
	// implementations.
	colReplace := strings.NewReplacer(
		TagDecorate(TextFalse, "td")+"</tr>", "</tr>",
		TagDecorate(TextTrue, "td")+"</tr>", "</tr>")

	var requestTableRows, responseTableRows string
	if FlagDebug {
		xLog.Printf("api is type %T", api)
	}

	for _, val01 := range api.Paths {
		requestTableRows = ""
		// avoid issues if no application/json body
		if schema, ok := val01.Post.RequestBody.Value.Content[KeyAppJson]; ok {
			requestTableRows = getSchemaProperties(schema.Schema, "", nil)
		}

		responseTableRows = ""
		// avoid issues if no '200' success response
		if testResponse, ok := val01.Post.Responses["200"]; ok {
			// avoid issues if no application/json body
			if schema, ok := testResponse.Value.Content[KeyAppJson]; ok {
				responseTableRows = getSchemaProperties(schema.Schema, "", nil)
			}
		}
		/*
			if FlagDebug {
				xLog.Printf("api is type %T, %T", requestTableRows, responseTableRows)
			}
		*/

		var sb strings.Builder
		// if there is no data for the table, don't display a table
		if "" != requestTableRows {
			misc.WriteSB(&sb, RequestHeader, SimplexRequestTableHeader,
				TableBodyOpen, requestTableRows, TableBodyAndTableClose)
		}
		if "" != responseTableRows {
			misc.WriteSB(&sb, ResponseHeader, SimplexResponseTableHeader,
				TableBodyOpen, colReplace.Replace(responseTableRows), TableBodyAndTableClose)
		}
		if sb.Len() > 0 {
			val01.Post.Description =
				val01.Post.Description + sb.String()
		}
	}
}

/*
func getSchemaProperties(i *oas.Schema) (properties string)
fetch properties as a list of table rows
*/
func getSchemaProperties(j *oas.SchemaRef, paramName string, required []string) (properties string) {
	var sb strings.Builder
	var pivot string

	// for technical reasons, need to treat "" as "object"
	// and switch won't accept a "" case.
	if "" != j.Value.Type {
		pivot = j.Value.Type
	} else {
		pivot = "object"
	}

	switch pivot {
	case "boolean":
		fallthrough
	case "integer":
		fallthrough
	case "number":
		fallthrough
	case "string":
		return MakeTableRow(paramName, j.Value.Description, required)

	case "object":
		for key, val := range j.Value.Properties {
			sb.WriteString(getSchemaProperties(val, key, append(required, j.Value.Required...)))
		}
		return sb.String()

	case "array":
		return getSchemaProperties(j.Value.Items, paramName, append(required, j.Value.Required...))

	default:
		runtime.Breakpoint()
		xLog.Panicf("huh? Parameter %s has unrecognized parameter type %s",
			paramName, j.Value.Type)
	}
	// unreached code
	//xLog.Panic("No possible execution path for this statement!")
	return ""
}

// MakeTableRow turn a name description and required status to a table entry
func MakeTableRow(paramName string, description string, required []string) (tableRow string) {

	var requiredText = TextFalse
	var sb strings.Builder

	if "" == paramName {
		return ""
	}
	if "" == description {
		description = TagDecorate("no description", "b")
	}

	for _, v := range required {
		if strings.EqualFold(paramName, v) {
			requiredText = TextTrue
			// test for mismatched case error between
			// required and the parameter definition
			if paramName != v {
				sb.Reset()
				misc.WriteSB(&sb,
					"Warning: mismatch between parameter ",
					paramName, " and ", v, "\n",
					"The list of all required strings are: \n")

				xLog.Println(sb.String(), required)
				//fmt.Println(sb.String())
			}
			break
		}
	}

	// table data
	sb.Reset()
	misc.WriteSB(&sb,
		TagDecorate(paramName, "td"),
		TagDecorate(description, "td"),
		TagDecorate(requiredText, "td"))

	return TagDecorate(sb.String(), "tr")

}
