package lisp

import (
	"errors"
	"fmt"
)

type CEnv = map[string]*Location

type Compiler struct {
	insns []Insn
	cenv  CEnv
	level int
}

func NewCompiler() *Compiler {
	return &Compiler{cenv: CEnv{}}
}

func (c *Compiler) clone() *Compiler {
	cenv := CEnv{}
	for k, v := range c.cenv {
		cenv[k] = v
	}
	return &Compiler{nil, cenv, c.level}
}

func (c *Compiler) pushInsn(op Op, operands []Operand) {
	c.insns = append(c.insns, Insn{op, operands})
}

func (c *Compiler) compile(expr Object) error {
	switch e := expr.(type) {
	case nil:
		c.pushInsn(NIL, nil)
	case bool, int:
		c.pushInsn(LDC, []Operand{e})
	case *Symbol:
		loc := c.cenv[e.name]
		if loc == nil {
			return fmt.Errorf("unknown variable: %s", e.name)
		}
		c.pushInsn(LD, []Operand{&Location{c.level - loc.level, loc.offset}})
	case *Cons:
		return c.compileList(e.car, e.cdr)
	}
	return nil
}

func (c *Compiler) compileList(car Object, cdr Object) error {
	switch obj := car.(type) {
	case *Symbol:
		switch obj.name {
		case "+":
			return c.compileOp(2, cdr, ADD)
		case "-":
			return c.compileOp(2, cdr, SUB)
		case "*":
			return c.compileOp(2, cdr, MUL)
		case "/":
			return c.compileOp(2, cdr, DIV)
		case "=":
			return c.compileOp(2, cdr, EQ)
		case "<":
			return c.compileOp(2, cdr, LT)
		case ">":
			return c.compileOp(2, cdr, GT)
		case "<=":
			return c.compileOp(2, cdr, LTE)
		case ">=":
			return c.compileOp(2, cdr, GTE)
		case "cons":
			return c.compileOp(2, cdr, CONS)
		case "car":
			return c.compileOp(1, cdr, CAR)
		case "cdr":
			return c.compileOp(1, cdr, CDR)
		case "null":
			return c.compileOp(1, cdr, NULL)
		case "atom":
			return c.compileOp(1, cdr, ATOM)
		case "quote":
			return c.compileQuote(cdr)
		case "if":
			return c.compileIf(cdr)
		case "set!":
			return c.compileSet(cdr)
		case "begin":
			return c.compileBegin(cdr)
		case "lambda":
			return c.compileLambda(cdr)
		default:
			return c.compileApplication(car, cdr)
		}
	case *Cons:
		return c.compileApplication(car, cdr)
	default:
		return fmt.Errorf("%s is not applicable", ToString(car))
	}
}

func (c *Compiler) takeArgs(n int, argList Object) ([]Object, error) {
	ret, improper, err := ListToSlice(argList)
	if improper != nil || err != nil {
		return nil, errors.New("arglist must be proper list")
	}
	nargs := len(ret)
	if nargs < n {
		return nil, errors.New("too less arguments")
	} else if nargs > n {
		return nil, errors.New("too many arguments")
	}
	return ret, nil
}

func (c *Compiler) compileOp(nargs int, argList Object, op Op) error {
	args, err := c.takeArgs(nargs, argList)
	if err != nil {
		return err
	}
	for _, arg := range args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}
	c.pushInsn(op, nil)
	return nil
}

func (c *Compiler) compileQuote(argList Object) error {
	args, err := c.takeArgs(1, argList)
	if err != nil {
		return err
	}
	c.pushInsn(LDC, []Operand{args[0]})
	return nil
}

func (c *Compiler) compileIf(argList Object) error {
	args, err := c.takeArgs(3, argList)
	if err != nil {
		return err
	}
	c1 := c.clone()
	c2 := c.clone()
	if err := c.compile(args[0]); err != nil {
		return err
	}
	if err := c1.compile(args[1]); err != nil {
		return err
	}
	c1.pushInsn(JOIN, nil)
	if err := c2.compile(args[2]); err != nil {
		return err
	}
	c2.pushInsn(JOIN, nil)
	c.pushInsn(SEL, []Operand{Code(c1.insns), Code(c2.insns)})
	return nil
}

func (c *Compiler) compileSet(argList Object) error {
	args, err := c.takeArgs(2, argList)
	if err != nil {
		return err
	}
	binding, ok := args[0].(*Symbol)
	if !ok {
		return errors.New("first argument of set! must be a symbol")
	}
	loc := c.cenv[binding.name]
	if loc == nil {
		return fmt.Errorf("unknown variable: %s", binding.name)
	}
	if err = c.compile(args[1]); err != nil {
		return err
	}
	c.pushInsn(SV, []Operand{&Location{c.level - loc.level, loc.offset}})
	return nil
}

func (c *Compiler) compileExprs(exprs []Object) error {
	for i, expr := range exprs {
		if err := c.compile(expr); err != nil {
			return err
		}
		if i < len(exprs)-1 {
			c.pushInsn(POP, nil)
		}
	}
	return nil
}

func (c *Compiler) compileBegin(argList Object) error {
	exprs, improper, err := ListToSlice(argList)
	if improper != nil || err != nil {
		return errors.New("arglist must be proper list")
	}
	return c.compileExprs(exprs)
}

func (c *Compiler) compileLambda(argList Object) error {
	args, improper, err := ListToSlice(argList)
	if improper != nil || err != nil {
		return errors.New("arglist must be proper list")
	}
	cbody := c.clone()
	cbody.level++
	params, improper, err := ListToSlice(args[0])
	if improper != nil || err != nil {
		return errors.New("arglist must be proper list")
	}
	for i, param := range params {
		switch obj := param.(type) {
		case *Symbol:
			cbody.cenv[obj.name] = &Location{cbody.level, i}
		default:
			return errors.New("fn argument must be symbol")
		}
	}
	if err := cbody.compileExprs(args[1:]); err != nil {
		return err
	}
	cbody.pushInsn(RTN, nil)
	c.pushInsn(LDF, []Operand{Code(cbody.insns)})
	return nil
}

func (c *Compiler) compileApplication(fn Object, argList Object) error {
	args, improper, err := ListToSlice(argList)
	if improper != nil || err != nil {
		return errors.New("arglist must be proper list")
	}
	for _, arg := range args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}
	c.pushInsn(NIL, nil)
	for range args {
		c.pushInsn(CONS, nil)
	}
	if err := c.compile(fn); err != nil {
		return err
	}
	c.pushInsn(AP, nil)
	return nil
}

func Compile(expr Object) (Code, error) {
	compiler := NewCompiler()
	if err := compiler.compile(expr); err != nil {
		return nil, err
	}
	return compiler.insns, nil
}
