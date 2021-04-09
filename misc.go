package main

import (
	"fmt"
	"strings"

	oas "github.com/getkin/kin-openapi/openapi3"
)

// DeferError
// account for an at-close function that
// may return an error for its close
func DeferError(f func() error) {
	err := f()
	if nil != err {
		xLog.Printf("%s%s",
			"(may be harmless) error in deferred function: ",
			err.Error())
	}
}

// TagDecorate Add HTML open/close tag to a phrase, and
// accept any attributes as well such as:
// TagDecorate("Huh?", "b", "style=\"color: hotpink;\"")
// would produce:
// <b style="color: hotpink">Huh?</b>
func TagDecorate(base string, tags ...string) (decoratedString string) {
	var sb strings.Builder
	if 0 >= len(tags) {
		xLog.Printf("huh? Attempting to HTML-ize %s without an HTML tag?", base)
	}
	WriteSB(&sb, "<", tags[0])
	for ix01 := 1; ix01 < len(tags); ix01++ {
		if "" != tags[ix01] {
			WriteSB(&sb, " ", tags[ix01])
		}
	}
	WriteSB(&sb, ">", MarkdownToHtml(base), "</", tags[0], ">")
	return sb.String()
}

// WriteSB Add a series of strings to a strings.Builder
func WriteSB(sb *strings.Builder, input ...string) {
	for _, val := range input {
		_, err := sb.WriteString(val)
		if nil != err {
			xLog.Print("strings.Builder failed to add " + val + " ??")
			xLog.Fatal("values: ", input)
		}
	}
}

// MarkdownToHtml convert some simple markdown formatting
// into HTML-tagged formatted text, naively.
func MarkdownToHtml(markdown string) (html string) {
	const ThreeStars = "***"
	const TwoStars = "**"
	const OneStar = "*"
	const ThreeUnder = "___"
	const TwoUnder = "__"
	const OneUnder = "_"
	const Sneech = "`" // https://bit.ly/3qSumos

	if "" == markdown {
		return ""
	}
	// MIGHT be a better way to do it?  MUST use this order!
	html = MarkdownSeparatorToHtmlTags(markdown, ThreeStars, "b", "i")
	html = MarkdownSeparatorToHtmlTags(html, TwoStars, "b")
	html = MarkdownSeparatorToHtmlTags(html, OneStar, "i")
	html = MarkdownSeparatorToHtmlTags(html, ThreeUnder, "i", "b")
	html = MarkdownSeparatorToHtmlTags(html, TwoUnder, "b")
	html = MarkdownSeparatorToHtmlTags(html, OneUnder, "i")
	return MarkdownSeparatorToHtmlTags(html, Sneech, "code")
}

// MarkdownSeparatorToHtmlTags do the actual transformation
// of A MARKDOWN FORMAT MODE to HTML-TAGGED TEXT
func MarkdownSeparatorToHtmlTags(target string, sep string, tags ...string) (rep string) {
	var replace strings.Builder
	substr := strings.SplitAfter(target, sep)

	if nil == substr || 0 >= len(substr) || 1 != len(substr)%2 {
		if FlagDebug {
			replace.WriteString(
				fmt.Sprintf("huh?\n\tDESCRIPTION: %s\n\nSEPARATOR: %s\n\t TAGS: %q\n\t",
					target, sep, tags))
			if nil == substr {
				replace.WriteString("-- no generated substrings --")
			} else {
				replace.WriteString(
					fmt.Sprintf("unexpected (%d is even) number of substrings:\n\t[ %q ]",
						len(substr)%2, substr))
			}
			xLog.Println(replace.String())
			replace.Reset()
		}
		return target
	}

	if 2 < len(substr) { // any separations happen?
		for ix01 := 0; ix01 < len(substr)-1; ix01 += 2 {
			markedText := strings.ReplaceAll(substr[ix01+1], sep, "")
			for _, val := range tags {
				markedText = TagDecorate(markedText, val)
			}
			WriteSB(&replace,
				strings.ReplaceAll(substr[ix01], sep, ""),
				markedText)
		}
		WriteSB(&replace, substr[len(substr)-1])
	} else {
		return target
	}
	return replace.String()
}

// StripReferences Remove the reference pointers so that
// UnMarshal expands everything.
func StripReferences(api *oas.Swagger) {
	//  Technically could do this for all operations
	//  Connect, Delete, Get, Head,
	//  Options, Patch, Put, and Trace
	//  but our file only uses Post
	for _, val01 := range api.Paths {
		val01.Post.RequestBody.Ref = ""
		for _, val02 := range val01.Post.RequestBody.Value.Content {
			StripReferencesSchema(val02.Schema)
		}
		for _, val03 := range val01.Post.Responses {
			val03.Ref = ""
			for _, val04 := range val03.Value.Content {
				StripReferencesSchema(val04.Schema)
			}
		}
	}
}

// StripReferencesSchema Eventually, whatever reference
// paths an OAS file has comes down to schema references,
// and this recursively clears those references
func StripReferencesSchema(schema *oas.SchemaRef) {
	// clean this reference, and look for sub-references
	// within the schemaBody
	schema.Ref = ""
	schemaBody := schema.Value
	if 0 != len(schemaBody.Extensions) {
		xLog.Print("WARNING: Extension Properties are NOT handled by this program")
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
