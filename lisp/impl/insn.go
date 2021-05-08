package lisp

type Op int

const (
	NIL Op = iota
	LDC
	LD
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
)

type Operand interface{}
type Insn struct {
	operator Op
	operands []Operand
}
type Code []Insn
