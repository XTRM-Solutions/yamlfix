package misc

import (
	"fmt"
	"os"
	"strings"
)

// DeferError
// account for an at-close function that
// may return an error for its close
func DeferError(f func() error) {
	err := f()
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr,
			"got error in DeferError, may be harmless\nerror: %s\n",
			err.Error())
	}
}

// WriteSB Add a series of strings to a strings.Builder
func WriteSB(sb *strings.Builder, inputStrings ...string) {
	if nil == sb || nil == inputStrings {
		panic("null pointer instead of *strings.Builder or inputStrings in misc.WriteSB()")
	}
	if len(inputStrings) <= 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Got 0-length array of strings in misc.WriteSB()")
		return
	}
	for _, val := range inputStrings {
		_, err := sb.WriteString(val)
		if nil != err {
			_, _ = fmt.Fprintf(os.Stderr,
				"Got error in misc.WriteSB() while writing strings.\nError: %s",
				err.Error())
			for ix, str := range inputStrings {
				_, _ = fmt.Fprintf(os.Stderr, "%05d: [ %s ]", ix, str)
			}
			panic("panic in misc.WriteSB")
		}
	}
}
