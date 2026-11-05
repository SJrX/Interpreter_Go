package evaluator

import (
	"testing"

	"github.com/SJrX/Interpreter_Go/lexer"
	"github.com/SJrX/Interpreter_Go/object"
	"github.com/SJrX/Interpreter_Go/parser"
)

// These tests run whole Monkey programs through the lexer, parser and
// evaluator together, exercising features in combination rather than one
// node type at a time.

// prelude defines list helpers in Monkey itself, built only from recursion,
// closures and the built-in functions.
const prelude = `
let map = fn(arr, f) {
  let iter = fn(arr, acc) {
    if (len(arr) == 0) { acc } else { iter(rest(arr), push(acc, f(first(arr)))) }
  };
  iter(arr, [])
};

let filter = fn(arr, keep) {
  let iter = fn(arr, acc) {
    if (len(arr) == 0) {
      acc
    } else {
      let x = first(arr);
      if (keep(x)) { iter(rest(arr), push(acc, x)) } else { iter(rest(arr), acc) }
    }
  };
  iter(arr, [])
};

let reduce = fn(arr, initial, f) {
  let iter = fn(arr, acc) {
    if (len(arr) == 0) { acc } else { iter(rest(arr), f(acc, first(arr))) }
  };
  iter(arr, initial)
};

let range = fn(from, to) {
  let iter = fn(i, acc) {
    if (i > to) { acc } else { iter(i + 1, push(acc, i)) }
  };
  iter(from, [])
};
`

func evalWholeProgram(t *testing.T, input string) object.Object {
	t.Helper()
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors for %q: %v", input, p.Errors())
	}
	return Eval(program, object.NewEnvironment())
}

func requireInspect(t *testing.T, input string, want string) {
	t.Helper()
	got := evalWholeProgram(t, input)
	if got == nil {
		t.Fatalf("program evaluated to Go nil, want %q", want)
	}
	if got.Inspect() != want {
		t.Fatalf("wrong result.\nwant=%s\ngot =%s", want, got.Inspect())
	}
}

func TestProgramsUsingHigherOrderFunctions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"map", `map([1, 2, 3], fn(x) { x * x })`, "[1, 4, 9]"},
		{"filter", `filter(range(1, 10), fn(x) { x / 2 * 2 == x })`, "[2, 4, 6, 8, 10]"},
		{"reduce sum", `reduce(range(1, 100), 0, fn(acc, x) { acc + x })`, "5050"},
		{"map over empty", `map([], fn(x) { x })`, "[]"},
		{"pipeline", `reduce(map(filter(range(1, 6), fn(x) { x > 2 }), fn(x) { x * 10 }), 0, fn(a, b) { a + b })`, "180"},
		{"functions in arrays", `let fs = [fn(x) { x + 1 }, fn(x) { x * 2 }]; reduce(fs, 5, fn(acc, f) { f(acc) })`, "12"},
		{"functions in hashes", `let ops = {"add": fn(a, b) { a + b }, "mul": fn(a, b) { a * b }}; ops["mul"](6, ops["add"](3, 4))`, "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireInspect(t, prelude+tt.input, tt.want)
		})
	}
}

func TestProgramsUsingClosures(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"adder captures its argument",
			`let makeAdder = fn(n) { fn(x) { x + n } }; let addTwo = makeAdder(2); let addTen = makeAdder(10); [addTwo(1), addTen(1)]`,
			"[3, 11]",
		},
		{
			"closure sees bindings at creation time, not later shadowing",
			`let x = 1; let getX = fn() { x }; let f = fn() { let x = 2; getX() }; f()`,
			"1",
		},
		{
			"inner let does not leak to the caller",
			`let x = 1; let f = fn() { let x = 99; x }; [f(), x]`,
			"[99, 1]",
		},
		{
			"currying three levels deep",
			`let addThree = fn(a) { fn(b) { fn(c) { a + b + c } } }; addThree(1)(2)(3)`,
			"6",
		},
		{
			"compose",
			`let compose = fn(f, g) { fn(x) { f(g(x)) } }; let inc = fn(x) { x + 1 }; let dbl = fn(x) { x * 2 }; [compose(inc, dbl)(5), compose(dbl, inc)(5)]`,
			"[11, 12]",
		},
		{
			"parameter shadows outer binding",
			`let n = 100; let f = fn(n) { n * 2 }; [f(3), n]`,
			"[6, 100]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireInspect(t, tt.input, tt.want)
		})
	}
}

func TestProgramsUsingRecursion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"fibonacci", `let fib = fn(n) { if (n < 2) { n } else { fib(n - 1) + fib(n - 2) } }; fib(20)`, "6765"},
		{"mutual recursion", `let isEven = fn(n) { if (n == 0) { true } else { isOdd(n - 1) } }; let isOdd = fn(n) { if (n == 0) { false } else { isEven(n - 1) } }; [isEven(10), isOdd(7), isEven(3)]`, "[true, true, false]"},
		{"ackermann", `let ack = fn(m, n) { if (m == 0) { n + 1 } else { if (n == 0) { ack(m - 1, 1) } else { ack(m - 1, ack(m, n - 1)) } } }; ack(2, 3)`, "9"},
		{"deep recursion", `let count = fn(n) { if (n == 0) { 0 } else { 1 + count(n - 1) } }; count(5000)`, "5000"},
		{"reverse", `let reverse = fn(arr) { if (len(arr) == 0) { [] } else { push(reverse(rest(arr)), first(arr)) } }; reverse([1, 2, 3, 4])`, "[4, 3, 2, 1]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireInspect(t, tt.input, tt.want)
		})
	}
}

func TestProgramsReturnFromNestedBlocks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"return from nested if", `let f = fn(x) { if (x > 0) { if (x > 10) { return "big"; } return "small"; } "non-positive" }; [f(50), f(5), f(-1)]`, `[big, small, non-positive]`},
		{"return stops the function, not the caller", `let f = fn() { return 1; 2 }; let g = fn() { f(); 3 }; g()`, "3"},
		{"statements after return are not evaluated", `let f = fn() { return 1; undefinedVariable }; f()`, "1"},
		{"top-level return stops the program", `return 7; 8`, "7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireInspect(t, tt.input, tt.want)
		})
	}
}

// An error anywhere must stop evaluation and be reported, wherever it occurs.
func TestProgramsPropagateErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"inside a function argument", `let f = fn(x) { x }; f(1 + true)`, "ERROR: type mismatch: INTEGER + BOOLEAN"},
		{"inside an if condition", `if (-true) { 1 } else { 2 }`, "ERROR: unknown operator: -BOOLEAN"},
		{"inside an array literal", `[1, 2, missing, 4]`, "ERROR: identifier not found: missing"},
		{"inside a hash value", `{"a": 1, "b": len(1)}`, "ERROR: argument to `len` not supported, got INTEGER"},
		{"inside an index", `[1, 2, 3][true + 1]`, "ERROR: type mismatch: BOOLEAN + INTEGER"},
		{"inside a deeply nested call", `let f = fn(n) { if (n == 0) { nope } else { f(n - 1) } }; f(10)`, "ERROR: identifier not found: nope"},
		{"stops the rest of the program", `let a = 1 + true; let b = missing; b`, "ERROR: type mismatch: INTEGER + BOOLEAN"},
		{"inside a higher-order callback", prelude + `map([1, 2, 3], fn(x) { x + "a" })`, "ERROR: type mismatch: INTEGER + STRING"},
		{"calling a non-function", `let x = 5; x(1)`, "ERROR: not a function: INTEGER"},
		{"unhashable key", `{[1]: 2}`, "ERROR: unusable as hash key: ARRAY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireInspect(t, tt.input, tt.want)
		})
	}
}

func TestProgramsTreatBuiltinsAsImmutable(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"push does not modify its argument", `let a = [1, 2]; let b = push(a, 3); [a, b]`, "[[1, 2], [1, 2, 3]]"},
		{"rest does not modify its argument", `let a = [1, 2, 3]; let b = rest(a); [a, b]`, "[[1, 2, 3], [2, 3]]"},
		{"rest of a single element", `rest([1])`, "[]"},
		{"builtins on empty arrays", `[first([]), last([]), rest([])]`, "[null, null, null]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireInspect(t, tt.input, tt.want)
		})
	}
}

func TestProgramsUsingHashes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"mixed key types", `let h = {"one": 1, 2: "two", true: "yes"}; [h["one"], h[2], h[true]]`, "[1, two, yes]"},
		{"keys computed from expressions", `let k = "na" + "me"; let h = {k: "Monkey", 1 + 1: "two"}; [h["name"], h[2]]`, "[Monkey, two]"},
		{"missing key", `{"a": 1}["b"]`, "null"},
		{"hash of arrays", `let people = {"names": ["Ada", "Alan"]}; people["names"][1]`, "Alan"},
		{"array of hashes", `let rows = [{"n": 1}, {"n": 2}]; rows[0]["n"] + rows[1]["n"]`, "3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireInspect(t, tt.input, tt.want)
		})
	}
}
