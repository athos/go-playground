package lisp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM(t *testing.T) {
	tests := []struct {
		title string
		code  Code
		out   interface{}
	}{
		{
			"nil -> nil",
			[]Insn{{NIL, nil}},
			nil,
		},
		{
			"ldc(t) -> t",
			[]Insn{{LDC, []Operand{true}}},
			true,
		},
		{
			"nil; null -> t",
			[]Insn{{NIL, nil}, {NULL, nil}},
			true,
		},
		{
			"nil; atom -> t",
			[]Insn{{NIL, nil}, {ATOM, nil}},
			true,
		},
		{
			"ldc(1); ldc(2); cons -> (1 . 2)",
			[]Insn{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{CONS, nil},
			},
			&Cons{1, 2},
		},
		{
			"ldc(1); ldc(2); cons; car -> 1",
			[]Insn{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{CONS, nil},
				{CAR, nil},
			},
			1,
		},
		{
			"ldc(1); ldc(2); cons; cdr -> 2",
			[]Insn{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{CONS, nil},
				{CDR, nil},
			},
			2,
		},
		{
			"ldc(1); ldc(2); cons; atom -> nil",
			[]Insn{
				{LDC, []Operand{1}},
				{LDC, []Operand{2}},
				{CONS, nil},
				{ATOM, nil},
			},
			nil,
		},
		{
			"ldc(2); ldc(1); add; -> 3",
			[]Insn{
				{LDC, []Operand{2}},
				{LDC, []Operand{1}},
				{ADD, nil},
			},
			3,
		},
		{
			"ldc(5); ldc(3); sub; -> 2",
			[]Insn{
				{LDC, []Operand{5}},
				{LDC, []Operand{3}},
				{SUB, nil},
			},
			2,
		},
		{
			"ldc(2); ldc(3); mul; -> 6",
			[]Insn{
				{LDC, []Operand{2}},
				{LDC, []Operand{3}},
				{MUL, nil},
			},
			6,
		},
		{
			"ldc(8); ldc(4); mul; -> 2",
			[]Insn{
				{LDC, []Operand{8}},
				{LDC, []Operand{4}},
				{DIV, nil},
			},
			2,
		},
		{
			"ldc(2); ldc(2); eq; -> t",
			[]Insn{
				{LDC, []Operand{2}},
				{LDC, []Operand{2}},
				{EQ, nil},
			},
			true,
		},
		{
			"ldc(2); ldc(3); eq; -> nil",
			[]Insn{
				{LDC, []Operand{2}},
				{LDC, []Operand{3}},
				{EQ, nil},
			},
			nil,
		},
		{
			"ldc(5); ldc(3); gt; -> t",
			[]Insn{
				{LDC, []Operand{5}},
				{LDC, []Operand{3}},
				{GT, nil},
			},
			true,
		},
		{
			"ldc(2); ldc(3); gt; -> nil",
			[]Insn{
				{LDC, []Operand{2}},
				{LDC, []Operand{3}},
				{GT, nil},
			},
			nil,
		},
		{
			"ldc(2); ldc(3); lt; -> t",
			[]Insn{
				{LDC, []Operand{2}},
				{LDC, []Operand{3}},
				{LT, nil},
			},
			true,
		},
		{
			"ldc(5); ldc(3); lt; -> t",
			[]Insn{
				{LDC, []Operand{5}},
				{LDC, []Operand{3}},
				{LT, nil},
			},
			nil,
		},
		{
			"ldc(5); ldc(3); gte; -> t",
			[]Insn{
				{LDC, []Operand{5}},
				{LDC, []Operand{3}},
				{GTE, nil},
			},
			true,
		},
		{
			"ldc(2); ldc(3); gte; -> nil",
			[]Insn{
				{LDC, []Operand{2}},
				{LDC, []Operand{3}},
				{GTE, nil},
			},
			nil,
		},
		{
			"ldc(2); ldc(3); lte; -> t",
			[]Insn{
				{LDC, []Operand{2}},
				{LDC, []Operand{3}},
				{LTE, nil},
			},
			true,
		},
		{
			"ldc(5); ldc(3); lte; -> t",
			[]Insn{
				{LDC, []Operand{5}},
				{LDC, []Operand{3}},
				{LTE, nil},
			},
			nil,
		},
		{
			"nil; null; sel(ldc(1); join, ldc(2); join); -> 1",
			[]Insn{
				{NIL, nil},
				{NULL, nil},
				{SEL, []Operand{
					Code([]Insn{{LDC, []Operand{1}}, {JOIN, nil}}),
					Code([]Insn{{LDC, []Operand{2}}, {JOIN, nil}}),
				}},
			},
			1,
		},
		{
			"ldc(t); null; sel(ldc(1); join, ldc(2); join); -> 2",
			[]Insn{
				{LDC, []Operand{true}},
				{NULL, nil},
				{SEL, []Operand{
					Code([]Insn{{LDC, []Operand{1}}, {JOIN, nil}}),
					Code([]Insn{{LDC, []Operand{2}}, {JOIN, nil}}),
				}},
			},
			2,
		},
		{
			// ((lambda (x) (+ x 2)) 3)
			"ldc(3); nil; cons; ldf(ld(0,0); ldc(2); mul; rtn;); ap; -> 6",
			[]Insn{
				{LDC, []Operand{3}},
				{NIL, nil},
				{CONS, nil},
				{LDF, []Operand{
					Code([]Insn{
						{LD, []Operand{&Location{0, 0}}},
						{LDC, []Operand{2}},
						{MUL, nil},
						{RTN, nil},
					}),
				}},
				{AP, nil},
			},
			6,
		},
		{
			// (let ((f (lambda (x)
			//            (* x 2))))
			//   (f (f 3)))
			"ldf(ld(0,0); ldc(2); mul; rtn;); nil; cons; cons; ldf(ldc(3); nil; cons; ld(0,0); ap; nil; cons; ld(0,0); ap; rtn;); ap; -> 12",
			[]Insn{
				{LDF, []Operand{
					Code{
						{LD, []Operand{&Location{0, 0}}},
						{LDC, []Operand{2}},
						{MUL, nil},
						{RTN, nil},
					},
				}},
				{NIL, nil},
				{CONS, nil},
				{LDF, []Operand{
					Code{
						{LDC, []Operand{3}},
						{NIL, nil},
						{CONS, nil},
						{LD, []Operand{&Location{0, 0}}},
						{AP, nil},
						{NIL, nil},
						{CONS, nil},
						{LD, []Operand{&Location{0, 0}}},
						{AP, nil},
						{RTN, nil},
					},
				}},
				{AP, nil},
			},
			12,
		},
		{
			// (letrec ((f (lambda (x)
			//               (if (= x 0)
			//                 1
			//                 (* x (f (- x 1)))))))
			//   (f 5))
			"ldc(5); nil; cons; ldf(ld(0,0); ldc(0); eq; sel(ldc(1); join;, ld(0,0); ld(0,0); ldc(1); sub; nil; cons; ld(1,0); dum; rap; mul; join;); rtn;); dum; rap; -> 120",
			[]Insn{
				{LDC, []Operand{5}},
				{NIL, nil},
				{CONS, nil},
				{LDF, []Operand{
					Code{
						{LD, []Operand{&Location{0, 0}}},
						{LDC, []Operand{0}},
						{EQ, nil},
						{SEL, []Operand{
							Code{{LDC, []Operand{1}}, {JOIN, nil}},
							Code{
								{LD, []Operand{&Location{0, 0}}},
								{LD, []Operand{&Location{0, 0}}},
								{LDC, []Operand{1}},
								{SUB, nil},
								{NIL, nil},
								{CONS, nil},
								{LD, []Operand{&Location{1, 0}}},
								{DUM, nil},
								{RAP, nil},
								{MUL, nil},
								{JOIN, nil},
							},
						}},
						{RTN, nil},
					},
				}},
				{DUM, nil},
				{RAP, nil},
			},
			120,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			vm := NewVM(tt.code)
			v, err := vm.Run()
			assert.Equal(t, tt.out, v)
			assert.Equal(t, nil, err)
		})
	}
}
