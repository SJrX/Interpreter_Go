// Parser tracing helpers, for debugging the Pratt parser.
//
// They are not called anywhere by default. To see how an expression is
// parsed, add a call such as
//
//	defer untrace(trace("parseExpression"))
//
// at the top of a parse function and run the parser tests with -v.

//lint:file-ignore U1000 Debugging helpers, only used while investigating the parser.

package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
