# Monkey interpreter in Go

A tree-walking interpreter for **Monkey**, the small programming language
designed in Thorsten Ball's book
[*Writing An Interpreter In Go*](https://interpreterbook.com/).

The source code is lexed into tokens, parsed into an abstract syntax tree (AST)
with a Pratt parser, and evaluated by walking that tree directly. There is no
bytecode or virtual machine.

## Attribution

This repository is a fork of
[SaladinoBelisario/Interpreter_Go](https://github.com/SaladinoBelisario/Interpreter_Go)
by Saladino Belisario, who wrote the interpreter while working through the book.
That code is MIT-licensed. See [LICENSE](LICENSE).

The Monkey language, and the design of this interpreter, come from
[*Writing An Interpreter In Go*](https://interpreterbook.com/) by
Thorsten Ball. If you find this project useful, the book explains every part of
it step by step and is well worth reading. This project is not affiliated with
the book or its author.

Changes made in this fork are listed in the Git history.

## The language

Monkey has integers, booleans, strings, arrays, hashes, first-class functions
and closures, with a C-like syntax.

```
let total = (2 + 3) * 4;                      // 20
let greeting = "Hello, " + "Monkey";          // "Hello, Monkey"

let colours = ["red", "green", "blue"];
colours[1];                                   // "green"

let book = {"title": "Monkey", "pages": 42};
book["pages"];                                // 42
```

Functions are values. The last expression in a function body is its result,
and `return` is also available:

```
let square = fn(n) { n * n };
square(7);                                    // 49

let factorial = fn(n) {
  if (n < 2) { 1 } else { n * factorial(n - 1) }
};
factorial(5);                                 // 120
```

Functions can take and return other functions, and closures capture their
environment:

```
let applyTwice = fn(f, x) { f(f(x)) };
applyTwice(fn(x) { x * 3 }, 2);               // 18

let makeAdder = fn(n) { fn(x) { x + n } };
let addTen = makeAdder(10);
addTen(5);                                    // 15
```

### Built-in functions

| Function | Description | Example |
|---|---|---|
| `len(x)` | Length of a string or array | `len("monkey")` → `6` |
| `first(arr)` | First element of an array | `first([1, 2, 3])` → `1` |
| `last(arr)` | Last element of an array | `last([1, 2, 3])` → `3` |
| `rest(arr)` | A new array without the first element | `rest([1, 2, 3])` → `[2, 3]` |
| `push(arr, x)` | A new array with `x` appended | `push([1, 2], 3)` → `[1, 2, 3]` |
| `puts(x...)` | Print values, returns `null` | `puts("hi")` |

### Operators

- Arithmetic on integers: `+`, `-`, `*`, `/`
- Comparison: `<`, `>`, `==`, `!=`
- Prefix: `-` (negation), `!` (logical not)
- String concatenation: `+`
- Indexing: `array[i]`, `hash[key]`

### Known limitations

These come from the book's design. They are not bugs:

- Integers only. There are no floating-point numbers.
- No `<=` or `>=`.
- Identifiers may contain only letters and `_` (no digits).
- No loops. Use recursion instead.
- No comments in Monkey source.

## Getting started

Requires Go 1.22 or newer.

Run the REPL:

```
go run .
```

```
>> let makeAdder = fn(n) { fn(x) { x + n } };
>> let addTen = makeAdder(10);
>> addTen(5)
15
```

Run the tests:

```
go test ./...
```

## Project layout

| Package | Responsibility |
|---|---|
| `token` | Token types produced by the lexer |
| `lexer` | Turns source text into tokens |
| `ast` | The abstract syntax tree node types |
| `parser` | Pratt parser that builds the AST from tokens |
| `object` | Runtime values (integers, strings, functions, …) and environments |
| `evaluator` | Walks the AST and evaluates it, plus the built-in functions |
| `repl` | The interactive read-eval-print loop |

## License

MIT. See [LICENSE](LICENSE). The original copyright notice from
SaladinoBelisario/Interpreter_Go is kept there unchanged.
