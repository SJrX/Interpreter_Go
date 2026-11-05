package ast

import (
	"testing"

	"github.com/SJrX/Interpreter_Go/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestStringWithMissingExpressions(t *testing.T) {
	one := &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1}
	ident := &Identifier{Token: token.Token{Type: token.IDENT, Literal: "f"}, Value: "f"}

	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{"prefix", &PrefixExpression{Operator: "-"}, "(-)"},
		{"infix right", &InfixExpression{Left: one, Operator: "+"}, "(1 + )"},
		{"infix left", &InfixExpression{Operator: "+", Right: one}, "( + 1)"},
		{"call argument", &CallExpression{Function: ident, Arguments: []Expression{one, nil}}, "f(1, )"},
		{"call function", &CallExpression{Arguments: []Expression{one}}, "(1)"},
		{"array element", &ArrayLiteral{Elements: []Expression{nil, one}}, "[, 1]"},
		{"index", &IndexExpression{Left: ident}, "(f[])"},
		{"if condition", &IfExpression{Consequence: &BlockStatement{}}, "if "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.String(); got != tt.expected {
				t.Fatalf("want %q, got %q", tt.expected, got)
			}
		})
	}
}
