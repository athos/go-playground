package lisp

type Op int

const (
	NIL Op = iota
	LDC
	LD
	LDG
	SV
	SVG
	POP
	ATOM
	NULL
	CAR
	CDR
	CONS
	ADD
	SUB
	MUL
	DIV
	EQ
	GT
	LT
	GTE
	LTE
	SEL
	JOIN
	LDF
	AP
	RTN
	DUM
	RAP
)

type Operand interface{}
type Insn struct {
	operator Op
	operands []Operand
}
type Code []Insn
