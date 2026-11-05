package evaluator

import (
	"testing"

	"github.com/SJrX/Interpreter_Go/object"
)

func TestStringConversion(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{`str(42)`, "42"}, {`str(true)`, "true"}, {`str(false)`, "false"},
		{`str([1, 2])`, "[1, 2]"}, {`str("hello")`, "hello"}, {`str("")`, ""},
		{`"n = " + str(5)`, "n = 5"}, {`str(if (false) { 1 })`, "null"},
		{`str({"a": 1})`, "{a: 1}"}, {`str(len)`, "builtin function"},
		{`str(fn(x) { x })`, "fn(x) {\nx\n}"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, ok := testEval(tc.input).(*object.String)
			if !ok {
				t.Fatalf("expected string from %s", tc.input)
			}
			if got.Value != tc.want {
				t.Errorf("got %q, want %q", got.Value, tc.want)
			}
		})
	}
}

func TestIntegerConversion(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  int64
	}{
		{`int("42")`, 42}, {`int("-7")`, -7}, {`int("+7")`, 7}, {`int("08")`, 8},
		{`int(5)`, 5}, {`int("0")`, 0}, {`int("9223372036854775807")`, 9223372036854775807},
		{`int("-9223372036854775808")`, -9223372036854775808},
		{`int(str(-7))`, -7}, {`int(str(int("-9223372036854775808")))`, -9223372036854775808},
	} {
		t.Run(tc.input, func(t *testing.T) { testIntegerObject(t, testEval(tc.input), tc.want) })
	}
	if !testBooleanObject(t, testEval(`int(str(42)) == 42`), true) {
		t.Fail()
	}
}

func TestConversionErrors(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{`int("abc")`, `could not convert "abc" to INTEGER`},
		{`int("")`, `could not convert "" to INTEGER`},
		{`int("1.2")`, `could not convert "1.2" to INTEGER`},
		{`int(" 42")`, `could not convert " 42" to INTEGER`},
		{`int("0x10")`, `could not convert "0x10" to INTEGER`},
		{`int("9223372036854775808")`, `could not convert "9223372036854775808" to INTEGER`},
		{`int("-9223372036854775809")`, `could not convert "-9223372036854775809" to INTEGER`},
		{`int(true)`, "argument to `int` not supported, got BOOLEAN"},
		{`int([])`, "argument to `int` not supported, got ARRAY"},
		{`int(if (false) { 1 })`, "argument to `int` not supported, got NULL"},
		{`int()`, "wrong number of arguments. got=0, want=1"},
		{`int(1, 2)`, "wrong number of arguments. got=2, want=1"},
		{`str()`, "wrong number of arguments. got=0, want=1"},
		{`str(1, 2)`, "wrong number of arguments. got=2, want=1"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, ok := testEval(tc.input).(*object.Error)
			if !ok {
				t.Fatalf("expected error from %s", tc.input)
			}
			if got.Message != tc.want {
				t.Errorf("got %q, want %q", got.Message, tc.want)
			}
		})
	}
}
