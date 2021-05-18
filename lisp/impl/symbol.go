package lisp

type Symbol struct {
	name string
	value Object
}

var symbolTable = map[string]*Symbol{}

func Intern(name string) *Symbol {
	sym, ok := symbolTable[name]
	if !ok {
		sym = &Symbol{name: name}
		symbolTable[name] = sym
	}
	return sym
}

func (sym *Symbol) SetValue(val Object) {
	sym.value = val
}
