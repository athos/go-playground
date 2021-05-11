package lisp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		in  Object
		out Code
	}{
		{nil, Code{{NIL, nil}}},
		{true, Code{{LDC, []Operand{true}}}},
		{42, Code{{LDC, []Operand{42}}}},
		{
			&Cons{&Symbol{"+"}, &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{ADD, nil},
			},
		},
		{
			&Cons{&Symbol{"-"}, &Cons{5, &Cons{3, nil}}},
			Code{
				{LDC, []Operand{5}},
				{LDC, []Operand{3}},
				{SUB, nil},
			},
		},
		{
			&Cons{&Symbol{"*"}, &Cons{2, &Cons{3, nil}}},
			Code{
				{LDC, []Operand{2}},
				{LDC, []Operand{3}},
				{MUL, nil},
			},
		},
		{
			&Cons{&Symbol{"/"}, &Cons{8, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{8}},
				{LDC, []Operand{2}},
				{DIV, nil},
			},
		},
		{
			&Cons{&Symbol{"="}, &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{EQ, nil},
			},
		},
		{
			&Cons{&Symbol{"<"}, &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{LT, nil},
			},
		},
		{
			&Cons{&Symbol{">"}, &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{GT, nil},
			},
		},
		{
			&Cons{&Symbol{"<="}, &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{LTE, nil},
			},
		},
		{
			&Cons{&Symbol{">="}, &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{GTE, nil},
			},
		},
		{
			&Cons{&Symbol{"null"}, &Cons{nil, nil}},
			Code{
				{NIL, nil},
				{NULL, nil},
			},
		},
		{
			&Cons{&Symbol{"atom"}, &Cons{true, nil}},
			Code{
				{LDC, []Operand{true}},
				{ATOM, nil},
			},
		},
		{
			&Cons{&Symbol{"if"}, &Cons{true, &Cons{1, &Cons{2, nil}}}},
			Code{
				{LDC, []Operand{true}},
				{SEL, []Operand{
					Code{{LDC, []Operand{1}}, {JOIN, nil}},
					Code{{LDC, []Operand{2}}, {JOIN, nil}},
				}},
			},
		},
		{
			&Cons{
				&Cons{
					&Symbol{"lambda"},
					&Cons{
						&Cons{&Symbol{"x"}, nil},
						&Cons{
							&Cons{
								&Symbol{"*"},
								&Cons{&Symbol{"x"}, &Cons{3, nil}},
							},
							nil,
						},
					},
				},
				&Cons{2, nil},
			},
			Code{
				{NIL, nil},
				{LDC, []Operand{2}},
				{CONS, nil},
				{LDF, []Operand{
					Code{
						{LD, []Operand{&Location{0, 0}}},
						{LDC, []Operand{3}},
						{MUL, nil},
						{RTN, nil},
					},
				}},
				{AP, nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(ToString(tt.in), func(t *testing.T) {
			code, err := Compile(tt.in)
			assert.Equal(t, tt.out, code)
			assert.Nil(t, err)
		})
	}
}
