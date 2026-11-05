package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/SJrX/Interpreter_Go/ast"
)

type ObjectType string

const (
	NULL_OBJ         = "NULL"
	ERROR_OBJ        = "ERROR"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type Null struct{}

type ReturnValue struct {
	Value Object
}

type Error struct {
	Message string
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

type String struct {
	Value string
}

type Builtin struct {
	Fn BuiltinFunction
}

type BuiltinFunction func(args ...Object) Object

type Array struct {
	Elements []Object
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
	// Keys holds the keys of Pairs in insertion order, so that hashes print
	// in the order they were written.
	Keys []HashKey
}

// Set adds or replaces a pair, whose key must be Hashable. A new key goes at
// the end of the order; an existing key keeps its position and gets the new
// value.
func (h *Hash) Set(pair HashPair) {
	key := pair.Key.(Hashable).HashKey()
	if h.Pairs == nil {
		h.Pairs = make(map[HashKey]HashPair)
	}
	if _, exists := h.Pairs[key]; !exists {
		h.Keys = append(h.Keys, key)
	}
	h.Pairs[key] = pair
}

type Hashable interface {
	HashKey() HashKey
}

func (f *Function) Type() ObjectType     { return FUNCTION_OBJ }
func (b *Boolean) Type() ObjectType      { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string       { return fmt.Sprintf("%t", b.Value) }
func (n *Null) Type() ObjectType         { return NULL_OBJ }
func (n *Null) Inspect() string          { return "null" }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (e *Error) Type() ObjectType        { return ERROR_OBJ }
func (e *Error) Inspect() string         { return "ERROR: " + e.Message }
func (i *Integer) Type() ObjectType      { return INTEGER_OBJ }
func (i *Integer) Inspect() string       { return fmt.Sprintf("%d", i.Value) }
func (s *String) Type() ObjectType       { return STRING_OBJ }
func (s *String) Inspect() string        { return s.Value }
func (b *Builtin) Type() ObjectType      { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string       { return "builtin function" }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, key := range h.Keys {
		pair := h.Pairs[key]
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
