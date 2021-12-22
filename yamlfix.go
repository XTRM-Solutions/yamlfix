package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"

	oas "github.com/getkin/kin-openapi/openapi3"
)

func main() {
	InitLog()
	// LIFO order for defer
	defer DeferError(xLogFile.Close)
	defer DeferError(xLogBuffer.Flush)
	InitFlags()

	xApi, err := oas.NewLoader().LoadFromFile(GetFlagString("infile"))
	if nil == xApi || nil != err {
		if nil != err {
			xLog.Printf("failed to load %s because %s",
				GetFlagString("infile"), err.Error())
		}
		os.Exit(-1)
	}

	if FlagDebug {
		writeJsonOASFile(xApi, "debug_pre_"+GetFlagString("outfile"))
	}

	if !FlagNoTables {
		format := GetFlagString("format")
		switch format {
		case "SIMPLEX":
			SimplexEnhanceDescriptions(xApi)
			break
		default:
			xLog.Printf("Huh? Somehow an unrecognized format [ %s ] was requested?", format)
		}

	}

	if FlagApiReport {
		ApiReport(xApi)
	}

	if FlagDereference {
		if FlagDebug {
			xLog.Print("Writing output file dereferenced, and unmodified file as debug_post_reference_")
			writeJsonOASFile(xApi, "debug_post_reference_"+GetFlagString("outfile"))
		} else {
			xLog.Print("Writing output file with internal referenced expanded")
		}
		StripReferences(xApi)
		writeJsonOASFile(xApi, GetFlagString("outfile"))
	} else {
		writeJsonOASFile(xApi, GetFlagString("outfile"))
		if FlagDebug {
			xLog.Print("Writing output file with references, and expanded file as debug_post_dereference_")
			StripReferences(xApi)
			writeJsonOASFile(xApi, "debug_post_dereference_"+GetFlagString("outfile"))
		} else {
			xLog.Print("writing output file with internal references")
		}
	}

}

func writeJsonOASFile(api *oas.T, fileName string) {
	output, err := api.MarshalJSON()
	if nil != err {
		xLog.Fatalf("Attempting to reconstruct API spec failed because: %s", err.Error())
	}

	src := bufio.NewReader(bytes.NewReader(output))
	decJson := json.NewDecoder(src)

	outFile, err := os.Create(fileName)
	if nil != err {
		xLog.Fatalf("Failed to create output JSON file %s because %s",
			fileName, err.Error())
	}
	defer DeferError(outFile.Close)

	bufferOutFile := bufio.NewWriter(outFile)
	defer DeferError(bufferOutFile.Flush)

	encJson := json.NewEncoder(bufferOutFile)
	encJson.SetIndent("", FlagIndentString)
	encJson.SetEscapeHTML(false)

	var m map[string]interface{}
	err = nil
	for ; nil == err; err = encJson.Encode(&m) {
		err = decJson.Decode(&m)
		if nil != err {
			break
		}
	}
	if io.EOF != err {
		xLog.Fatalf("Could not output JSON file %s because: %s", fileName, err.Error())
	}

}
