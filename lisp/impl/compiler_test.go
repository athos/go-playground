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
			&Cons{Intern("+"), &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{ADD, nil},
			},
		},
		{
			&Cons{Intern("-"), &Cons{5, &Cons{3, nil}}},
			Code{
				{LDC, []Operand{5}},
				{LDC, []Operand{3}},
				{SUB, nil},
			},
		},
		{
			&Cons{Intern("*"), &Cons{2, &Cons{3, nil}}},
			Code{
				{LDC, []Operand{2}},
				{LDC, []Operand{3}},
				{MUL, nil},
			},
		},
		{
			&Cons{Intern("/"), &Cons{8, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{8}},
				{LDC, []Operand{2}},
				{DIV, nil},
			},
		},
		{
			&Cons{Intern("="), &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{EQ, nil},
			},
		},
		{
			&Cons{Intern("<"), &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{LT, nil},
			},
		},
		{
			&Cons{Intern(">"), &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{GT, nil},
			},
		},
		{
			&Cons{Intern("<="), &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{LTE, nil},
			},
		},
		{
			&Cons{Intern(">="), &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{GTE, nil},
			},
		},
		{
			&Cons{Intern("cons"), &Cons{1, &Cons{2, nil}}},
			Code{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{CONS, nil},
			},
		},
		{
			&Cons{Intern("null"), &Cons{nil, nil}},
			Code{
				{NIL, nil},
				{NULL, nil},
			},
		},
		{
			&Cons{Intern("atom"), &Cons{true, nil}},
			Code{
				{LDC, []Operand{true}},
				{ATOM, nil},
			},
		},
		{
			&Cons{Intern("quote"), &Cons{Intern("foo"), nil}},
			Code{
				{LDC, []Operand{Intern("foo")}},
			},
		},
		{
			&Cons{Intern("if"), &Cons{true, &Cons{1, &Cons{2, nil}}}},
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
					Intern("lambda"),
					&Cons{
						&Cons{Intern("x"), nil},
						&Cons{
							&Cons{
								Intern("*"),
								&Cons{Intern("x"), &Cons{3, nil}},
							},
							nil,
						},
					},
				},
				&Cons{2, nil},
			},
			Code{
				{LDC, []Operand{2}},
				{NIL, nil},
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
		{
			&Cons{
				&Cons{
					Intern("lambda"),
					&Cons{
						&Cons{Intern("x"), nil},
						&Cons{
							&Cons{
								Intern("set!"),
								&Cons{
									Intern("x"),
									&Cons{
										&Cons{
											Intern("+"),
											&Cons{
												Intern("x"),
												&Cons{1, nil},
											},
										},
										nil,
									},
								},
							},
							&Cons{Intern("x"), nil},
						},
					},
				},
				&Cons{42, nil},
			},
			[]Insn{
				{LDC, []Operand{42}},
				{NIL, nil},
				{CONS, nil},
				{LDF, []Operand{
					Code{
						{LD, []Operand{&Location{0,0}}},
						{LDC, []Operand{1}},
						{ADD, nil},
						{SV, []Operand{&Location{0,0}}},
						{POP, nil},
						{LD, []Operand{&Location{0,0}}},
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
