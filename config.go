package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

var xLogFile *os.File
var xLog log.Logger

var FlagDeref bool
var FlagDebug bool
var FlagVerbose bool
var FlagQuiet bool
var nFlags *pflag.FlagSet

func InitLog() {
	var err error
	var logWriters []io.Writer

	xLogFile, err = os.OpenFile("yamlfix.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		xLog.Fatalf("error opening file: %v", err)
	}

	logWriters = append(logWriters, os.Stderr)
	logWriters = append(logWriters, xLogFile)

	out := io.MultiWriter(logWriters...)
	xLog.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	xLog.SetOutput(out)
}

func InitFlags() {

	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)

	nFlags.StringP("infile", "i", "input.yaml", "name of YAML file to process")
	nFlags.StringP("outfile", "o", "output.json", "name of processed (output) JSON file")
	// nFlags.IntP("indent", "n", 2, "Spaces to use for each indent level")
	nFlags.BoolP("help", "h", false, "Display help message and usage information")
	// nFlags.BoolP("tab", "t", false, "Use tab character for indent (sets indent to one character)")
	nFlags.BoolP("debug", "d", true,
		"Enable additional informational and operational logging output for debug purposes")
	nFlags.BoolP("quiet", "q", false, "Suppress superfluous output. Overrides verbose.")
	nFlags.BoolP("verbose", "v", true, "Supply informative messages")
	nFlags.StringP("format", "", "SIMPLEX", "Name of formatting module/method (currently 'SIMPLEX' only) ")
	nFlags.BoolP("expand-references", "x", true, "Expand internal and external references in POST methods")

	err := nFlags.Parse(os.Args[1:])

	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "\nerror parsing flags: %s\n", err.Error())
		os.Exit(-1)
	}

	FlagDebug = GetFlagBool("debug")

	FlagDeref = GetFlagBool("expand-references")

	if GetFlagBool("help") {
		_, thisCmd := filepath.Split(os.Args[0])
		_, _ = fmt.Fprint(os.Stdout, "\n", "usage for ", thisCmd, ":\n", nFlags.FlagUsages(), "\n")
		UsageMessage()
		os.Exit(0)
	}

	formatStr := strings.ToLower(GetFlagString("format"))
	if FlagDebug {
		_, _ = fmt.Fprintf(os.Stderr, "format style: got %s format mode\n", formatStr)
	}
	switch formatStr {
	case "simplex":
		break
	default:
		UsageMessage()
		_, _ = fmt.Fprintf(os.Stdout,
			"%s%s%s\n\t%s\n\t\t%s\n",
			"Bad/unknown format [ ",
			formatStr,
			" ] requested.",
			"Supported format(s) are:",
			"SIMPLEX:\t\tdefault formatting")
		os.Exit(-1)
	}

	FlagQuiet = GetFlagBool("quiet")
	if FlagQuiet {
		FlagVerbose = false
	} else {
		FlagVerbose = GetFlagBool("verbose")
	}

	if FlagVerbose {
		_, _ = fmt.Fprint(os.Stdout, "\nVerbose Mode Engaged\n")
	}
}

// UsageMessage /* UsageMessage
func UsageMessage() {
	_, _ = fmt.Println("\n\tInformative Usage Message Here")
}

// GetFlagBool /* GetFlagBool(key string) (value string)
func GetFlagBool(key string) (value bool) {
	var err error
	value, err = nFlags.GetBool(key)
	if nil != err {
		xLog.Printf("error fetching value for boolean flag [ %s ]\n", key)
		return false
	}
	return value
}

// GetFlagString /* GetFlagString(key string) (value string)
func GetFlagString(key string) (value string) {
	var err error
	value, err = nFlags.GetString(key)
	if nil != err {
		xLog.Printf("error fetching value for string flag [ %s ]\n", key)
		return ""
	}
	return value
}

/* GetFlagInt(key string) (value int)
Return a flag with string value
*/
/*
func GetFlagInt(key string) (value int) {
	var err error
	value, err = nFlags.GetInt(key)
	if nil != err {
		xLog.Printf("error fetching value for integer flag [ %s ]\n", key)
		return 0
	}
	return value
}
*/
