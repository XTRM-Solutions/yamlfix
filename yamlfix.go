package main

import (
	oas "github.com/getkin/kin-openapi/openapi3"
	"os"
	"path/filepath"
	// https://pkg.go.dev/github.com/getkin/kin-openapi@v0.53.0/openapi3
)

func main() {
	InitLog()
	defer DeferError(xLogFile.Close)
	InitFlags()

	xApi, err := oas.NewSwaggerLoader().LoadSwaggerFromFile(GetFlagString("infile"))
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

	EnhanceDescriptions(xApi)

	writeJsonOASFile(xApi, GetFlagString("outfile"))

}

func writeJsonOASFile(api *oas.Swagger, fileName string) {

	output, err := api.MarshalJSON()
	if nil != err {
		xLog.Fatalf("Attempting to reconstruct API spec failed because: %s", err.Error())
	}

	outFile, err := os.Create(fileName)
	if nil != err {
		xLog.Fatalf("Failed to create file %s because %s",
			fileName, err.Error())
	}
	defer DeferError(outFile.Close)

	byteCount, err := outFile.Write(output)
	if nil != err {
		xLog.Fatalf("Failed writing output file %s because: %s",
			fileName, err.Error())
	}
	if len(output) != byteCount {
		xLog.Fatalf("Only wrote %d bytes (of %d ) to %s",
			len(output), byteCount, fileName)
	}

	if FlagDebug || FlagVerbose {
		fileStat, err := outFile.Stat()
		if nil != err {
			xLog.Fatalf("Failed to read fileStat for %s", fileName)
		}
		filePath, err := filepath.Abs(fileStat.Name())
		if nil != err {
			xLog.Fatalf("Failed to get absolute path for %s", fileName)
		}
		xLog.Printf("Writing output to file %s", filePath)
	}

}
