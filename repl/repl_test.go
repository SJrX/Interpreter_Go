package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestStartEvaluatesEachLineInASharedEnvironment(t *testing.T) {
	in := strings.NewReader("let a = 5;\nlet b = a * 2;\na + b\n")
	var out bytes.Buffer

	Start(in, &out)

	want := PROMPT + PROMPT + PROMPT + "15\n" + PROMPT
	if out.String() != want {
		t.Fatalf("unexpected output.\nwant=%q\ngot =%q", want, out.String())
	}
}

func TestStartReportsParserErrorsAndContinues(t *testing.T) {
	in := strings.NewReader("let = 1;\n1 + 1\n")
	var out bytes.Buffer

	Start(in, &out)

	got := out.String()
	if !strings.Contains(got, "parser errors:") {
		t.Fatalf("expected parser errors in output, got %q", got)
	}
	if !strings.HasSuffix(got, "2\n"+PROMPT) {
		t.Fatalf("expected REPL to keep evaluating after an error, got %q", got)
	}
}
