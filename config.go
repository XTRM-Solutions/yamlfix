package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

var xLogFile *os.File
var xLogBuffer *bufio.Writer

var xLog log.Logger

/* var FlagPretty bool  */

var FlagIndentString string
var FlagDereference bool
var FlagDebug bool
var FlagVerbose bool
var FlagQuiet bool
var FlagNoTables bool
var FlagApiReport bool
var nFlags *pflag.FlagSet

func InitLog() {
	var err error
	var logWriters []io.Writer

	xLogFile, err = os.OpenFile("yamlfix.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		xLog.Fatalf("error opening file: %v", err)
	}
	xLogBuffer = bufio.NewWriter(xLogFile)
	logWriters = append(logWriters, os.Stderr)
	logWriters = append(logWriters, xLogBuffer)

	xLog.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	xLog.SetOutput(io.MultiWriter(logWriters...))
}

func InitFlags() {

	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)
	nFlags.BoolP("expand-only", "", false, "Do not create JSON parameter tables (implies -x false)")
	nFlags.StringP("infile", "i", "input.old_yaml", "name of YAML file to process")
	nFlags.StringP("outfile", "o", "output.json", "name of processed (output) JSON file")
	nFlags.IntP("indent", "n", 2, "Spaces to use for each indent level")
	nFlags.BoolP("help", "h", false, "Display help message and usage information")
	nFlags.BoolP("tab", "t", false, "Use tab character for indent (sets indent to one character)")
	nFlags.BoolP("debug", "d", false,
		"Enable additional informational and operational logging output for debug purposes")
	nFlags.BoolP("quiet", "q", false, "Suppress output to stdout and stderr (output still goes to logfile)")
	nFlags.BoolP("verbose", "v", false, "Supply informative messages")
	nFlags.StringP("format", "", "SIMPLEX", "Name of formatting module/method (currently 'SIMPLEX' only) ")
	nFlags.BoolP("expand-references", "x", true, "Expand internal and external references in POST methods")
	nFlags.BoolP("api-report", "a", false, "Provide a CSV report on API names and parameters for all functions")
	// nFlags.BoolP("prettyprint", "p", true, "Pretty-print JSON output")

	err := nFlags.Parse(os.Args[1:])

	// do quietness setup first
	FlagQuiet = GetFlagBool("quiet")
	if FlagQuiet {
		xLog.SetOutput(xLogBuffer)
		// shut off error messages to stderr, only log them
	}

	FlagVerbose = GetFlagBool("verbose")
	if FlagVerbose {
		xLog.Print("Verbose mode engaged ... ")
	}

	if nil != err {
		xLog.Fatalf("\nerror parsing flags: %s\n%s %s\n%s\n\t%v\n",
			err.Error(), "common issue: 2 hyphens for long-form arguments,",
			"1 hyphen for short-form argument",
			"Program arguments are:",
			os.Args)
	}

	FlagDebug = GetFlagBool("debug")

	// FlagPretty = GetFlagBool("prettyprint")

	if len(os.Args) <= 1 || GetFlagBool("help") {
		_, thisCmd := filepath.Split(os.Args[0])
		_, _ = fmt.Fprint(os.Stdout, "\n", "usage for ", thisCmd, ":\n", nFlags.FlagUsages(), "\n")
		UsageMessage()
		os.Exit(0)
	}

	formatStr := strings.ToLower(GetFlagString("format"))
	if FlagDebug {
		xLog.Printf("format style: got %s format mode\n", formatStr)
	}
	switch formatStr {
	case "simplex":
		break
	default:
		UsageMessage()
		xLog.Fatalf(
			"%s%s%s\n\t%s\n\t\t%s\n",
			"Bad/unknown format [ ",
			formatStr,
			" ] requested.",
			"Supported format(s) are:",
			"SIMPLEX:\t\tdefault formatting")
	}

	if GetFlagBool("tab") {
		FlagIndentString = "\t"
		if FlagVerbose {
			xLog.Print("Using single tab for indent, indent flag ignored")
		}
	} else {
		ind := GetFlagInt("indent")
		if MaxIndent < ind {
			xLog.Printf("Huh? Indent set too big (%d) resetting to MaxIndent value %d", ind, MaxIndent)
			ind = MaxIndent
		} else {
			if MinIndent > ind {
				xLog.Printf("Huh? Indent set too small (%d) resetting to MinIndent value %d", ind, MinIndent)
				ind = MinIndent
			}
		}
		FlagIndentString = strings.Repeat(" ", ind)
	}

	if GetFlagBool("expand-only") {
		FlagNoTables = true
		FlagDereference = true
	} else {
		FlagNoTables = false
		FlagDereference = GetFlagBool("expand-references")
	}

	FlagApiReport = GetFlagBool("api-report")
}

// UsageMessage /* UsageMessage
func UsageMessage() {
	_, _ = fmt.Println("\n\tInsert Informative Usage Message Here")
}

// GetFlagBool /* GetFlagBool(key string) (value string)
func GetFlagBool(key string) (value bool) {
	var err error
	value, err = nFlags.GetBool(key)
	if nil != err {
		xLog.Fatalf("error fetching value for boolean flag [ %s ]\n", key)
		return false
	}
	return value
}

// GetFlagString /* GetFlagString(key string) (value string)
func GetFlagString(key string) (value string) {
	var err error
	value, err = nFlags.GetString(key)
	if nil != err {
		xLog.Fatalf("error fetching value for string flag [ %s ]\n", key)
		return ""
	}
	return value
}

// GetFlagInt fetch the value of integer flag
func GetFlagInt(key string) (value int) {
	var err error
	value, err = nFlags.GetInt(key)
	if nil != err {
		xLog.Fatalf("error fetching value for integer flag [ %s ]\n", key)
		return 0
	}
	return value
}
