package lisp

import "errors"

type PC int
type Stack []Object
type Restorer interface {
	restore(*VM)
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
	return &VM{code: code}
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

func (vm *VM) pop() Object {
	if len(vm.stack) == 0 {
		panic("stack underflow")
	}
	obj := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return obj
}

func (vm *VM) dumpPop() Restorer {
	size := len(vm.dump)
	if size == 0 {
		panic("dump stack underflow")
	}
	ret := vm.dump[size-1]
	vm.dump = vm.dump[:size-1]
	return ret
}

func (vm *VM) binaryOp(op func(int, int) Object) error {
	obj1 := vm.pop()
	obj2 := vm.pop()
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
			vm.push(vm.env.Locate(loc))
		case LDG:
			sym := insn.operands[0].(*Symbol)
			vm.push(sym.value)
		case SV:
			loc := insn.operands[0].(*Location)
			obj := vm.pop()
			vm.env.Update(loc, obj)
			vm.push(obj)
		case SVG:
			sym := insn.operands[0].(*Symbol)
			obj := vm.pop()
			sym.SetValue(obj)
			vm.push(obj)
		case POP:
			vm.pop()
		case ATOM:
			obj := vm.pop()
			vm.push(FromBool(IsAtom(obj)))
		case NULL:
			obj := vm.pop()
			vm.push(FromBool(IsNull(obj)))
		case CONS:
			y := vm.pop()
			x := vm.pop()
			vm.push(NewCons(x, y))
		case CAR:
			obj := vm.pop()
			car, err := Car(obj)
			if err != nil {
				return nil, err
			}
			vm.push(car)
		case CDR:
			obj := vm.pop()
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
			vm.runSel(ct, cf)
			continue
		case JOIN:
			entry := vm.dumpPop()
			_ = entry.(*SelDumpEntry)
			entry.restore(vm)
		case LDF:
			code := insn.operands[0].(Code)
			env := vm.env
			vm.push(NewFunc(code, env))
		case AP:
			if err := vm.runAp(); err != nil {
				return nil, err
			}
			continue
		case RTN:
			entry := vm.dumpPop()
			_ = entry.(*ApDumpEntry)
			entry.restore(vm)
		case DUM:
			frame := make([]Object, 1)
			vm.env = vm.env.Push(frame)
		case RAP:
			if err := vm.runRap(); err != nil {
				return nil, err
			}
			continue
		}
		vm.pc++
	}
	return vm.pop(), nil
}

func (entry *SelDumpEntry) restore(vm *VM) {
	vm.code = entry.code
	vm.pc = entry.pc
}

func (vm *VM) runSel(ct, cf Code) {
	var c Code
	obj := vm.pop()
	if ToBool(obj) {
		c = ct
	} else {
		c = cf
	}
	vm.code, c = c, vm.code
	vm.dump = append(vm.dump, &SelDumpEntry{c, vm.pc})
	vm.pc = 0
}

func (entry *ApDumpEntry) restore(vm *VM) {
	v := vm.pop()
	vm.stack = append(entry.stack, v)
	vm.env = entry.env
	vm.code = entry.code
	vm.pc = entry.pc
}

func (vm *VM) withFn(f func(*Func) Restorer) error {
	obj := vm.pop()
	fn, ok := obj.(*Func)
	if !ok {
		return errors.New("cannot apply object other than function")
	}
	args := vm.pop()
	frame, improper, err := ListToSlice(args)
	if err != nil {
		return err
	}
	if improper != nil {
		return errors.New("improper lists are not allowed for arg lists")
	}
	entry := f(fn)
	vm.stack = nil
	vm.env = vm.env.Push(frame)
	vm.code = fn.code
	vm.dump = append(vm.dump, entry)
	vm.pc = 0
	return nil
}

func (vm *VM) runAp() error {
	return vm.withFn(func(_ *Func) Restorer {
		return &ApDumpEntry{
			stack: vm.stack,
			env:   vm.env,
			code:  vm.code,
			pc:    vm.pc,
		}
	})
}

func (vm *VM) runRap() error {
	return vm.withFn(func(fn *Func) Restorer {
		vm.env.frame[0] = fn
		return &ApDumpEntry{
			stack: vm.stack,
			env:   vm.env.Pop(),
			code:  vm.code,
			pc:    vm.pc,
		}
	})
}
