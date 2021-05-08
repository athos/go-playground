package lisp

import "errors"

type PC int
type Stack []Object
type Restorer interface {
	restore(*VM) error
}
type Dump []Restorer

type VM struct {
	stack Stack
	env   *Env
	code  Code
	dump  Dump
	pc    PC
}

type SelDumpEntry struct {
	code Code
	pc   PC
}

type ApDumpEntry struct {
	stack Stack
	env   *Env
	code  Code
	pc    PC
}

func NewVM(code Code) *VM {
	return &VM{
		env:  NewEnv(),
		code: code,
	}
}

func (vm *VM) fetchInsn() (*Insn, bool) {
	if int(vm.pc) >= len(vm.code) {
		return nil, false
	}
	return &vm.code[vm.pc], true
}

func (vm *VM) push(obj Object) {
	vm.stack = append(vm.stack, obj)
}

func (vm *VM) pop() (Object, error) {
	if len(vm.stack) == 0 {
		return nil, errors.New("Stack underflow")
	}
	obj := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return obj, nil
}

func (vm *VM) dumpPop() (Restorer, error) {
	size := len(vm.dump)
	if size == 0 {
		return nil, errors.New("dump stack underflow")
	}
	ret := vm.dump[size-1]
	vm.dump = vm.dump[:size-1]
	return ret, nil
}

func (vm *VM) binaryOp(op func(int, int) Object) error {
	obj1, err := vm.pop()
	if err != nil {
		return err
	}
	obj2, err := vm.pop()
	if err != nil {
		return err
	}
	y, err := ToNumber(obj1)
	if err != nil {
		return err
	}
	x, err := ToNumber(obj2)
	if err != nil {
		return err
	}
	vm.push(op(x, y))
	return nil
}

func (vm *VM) arithOp(op func(int, int) int) error {
	return vm.binaryOp(func(x, y int) Object {
		return op(x, y)
	})
}

func (vm *VM) logicalOp(op func(int, int) bool) error {
	return vm.binaryOp(func(x, y int) Object {
		return FromBool(op(x, y))
	})
}

func (vm *VM) Run() (Object, error) {
	for {
		insn, ok := vm.fetchInsn()
		if !ok {
			break
		}
		switch insn.operator {
		case NIL:
			vm.push(nil)
		case LDC:
			vm.push(insn.operands[0])
		case LD:
			loc := insn.operands[0].(*Location)
			obj, err := vm.env.Locate(loc)
			if err != nil {
				return nil, err
			}
			vm.push(obj)
		case ATOM:
			obj, err := vm.pop()
			if err != nil {
				return nil, err
			}
			vm.push(FromBool(IsAtom(obj)))
		case NULL:
			obj, err := vm.pop()
			if err != nil {
				return nil, err
			}
			vm.push(FromBool(IsNull(obj)))
		case CONS:
			x, err := vm.pop()
			if err != nil {
				return nil, err
			}
			y, err := vm.pop()
			if err != nil {
				return nil, err
			}
			vm.push(ToCons(x, y))
		case CAR:
			obj, err := vm.pop()
			if err != nil {
				return nil, err
			}
			car, err := Car(obj)
			if err != nil {
				return nil, err
			}
			vm.push(car)
		case CDR:
			obj, err := vm.pop()
			if err != nil {
				return nil, err
			}
			cdr, err := Cdr(obj)
			if err != nil {
				return nil, err
			}
			vm.push(cdr)
		case ADD:
			vm.arithOp(func(x, y int) int { return x + y })
		case SUB:
			vm.arithOp(func(x, y int) int { return x - y })
		case MUL:
			vm.arithOp(func(x, y int) int { return x * y })
		case DIV:
			vm.arithOp(func(x, y int) int { return x / y })
		case EQ:
			vm.logicalOp(func(x, y int) bool { return x == y })
		case GT:
			vm.logicalOp(func(x, y int) bool { return x > y })
		case LT:
			vm.logicalOp(func(x, y int) bool { return x < y })
		case GTE:
			vm.logicalOp(func(x, y int) bool { return x >= y })
		case LTE:
			vm.logicalOp(func(x, y int) bool { return x <= y })
		case SEL:
			ct := insn.operands[0].(Code)
			cf := insn.operands[1].(Code)
			if err := vm.runSel(ct, cf); err != nil {
				return nil, err
			}
		case JOIN:
			if err := vm.runJoin(); err != nil {
				return nil, err
			}
		case LDF:
			code := insn.operands[0].(Code)
			env := vm.env
			vm.push(&Func{code, env})
		case AP:
			if err := vm.runAp(); err != nil {
				return nil, err
			}
		case RTN:
			if err := vm.runRtn(); err != nil {
				return nil, err
			}
		}
		vm.pc++
	}
	if ret, err := vm.pop(); err != nil {
		return nil, nil
	} else {
		return ret, nil
	}
}

func (entry *SelDumpEntry) restore(vm *VM) error {
	vm.code = entry.code
	vm.pc = entry.pc
	return nil
}

func (vm *VM) runSel(ct, cf Code) error {
	var c Code
	obj, err := vm.pop()
	if err != nil {
		return err
	}
	if ToBool(obj) {
		c = ct
	} else {
		c = cf
	}
	vm.code, c = c, vm.code
	vm.dump = append(vm.dump, &SelDumpEntry{c, vm.pc})
	return nil
}

func (vm *VM) runJoin() error {
	entry, err := vm.dumpPop()
	if err != nil {
		return nil
	}
	selEntry, ok := entry.(*SelDumpEntry)
	if !ok {
		return errors.New("run into incoherent dump entry (ap)")
	}
	selEntry.restore(vm)
	return nil
}

func (entry *ApDumpEntry) restore(vm *VM) error {
	v, err := vm.pop()
	if err != nil {
		return err
	}
	vm.stack = append(entry.stack, v)
	vm.env = entry.env
	vm.code = entry.code
	vm.pc = entry.pc
	return nil
}

func (vm *VM) runAp() error {
	obj, err := vm.pop()
	if err != nil {
		return nil
	}
	fn, ok := obj.(*Func)
	if !ok {
		return errors.New("cannot apply object other than function")
	}
	args, err := vm.pop()
	if err != nil {
		return nil
	}
	frame, err := ListToSlice(args)
	if err != nil {
		return nil
	}
	entry := &ApDumpEntry{
		stack: vm.stack,
		env:   vm.env,
		code:  vm.code,
		pc:    vm.pc,
	}
	vm.stack = nil
	vm.env = Push(vm.env, frame)
	vm.code = fn.code
	vm.dump = append(vm.dump, entry)
	vm.pc = 0
	return nil
}

func (vm *VM) runRtn() error {
	entry, err := vm.dumpPop()
	if err != nil {
		return err
	}
	apDumpEntry, ok := entry.(*ApDumpEntry)
	if !ok {
		return errors.New("run into incoherent dump entry (sel)")
	}
	apDumpEntry.restore(vm)
	return nil
}
