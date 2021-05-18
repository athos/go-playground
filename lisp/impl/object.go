package lisp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Object interface{}
type Cons struct {
	car Object
	cdr Object
}
type Func struct {
	code Code
	env  *Env
}

func IsAtom(obj Object) bool {
	_, ok := obj.(*Cons)
	return !ok
}

func IsNull(obj Object) bool {
	return obj == nil
}

func ToBool(obj Object) bool {
	return !IsNull(obj)
}

func ToNumber(obj Object) (int, error) {
	n, ok := obj.(int)
	if !ok {
		return 0, errors.New("cannot be converted to number")
	}
	return n, nil
}

func Car(obj Object) (Object, error) {
	c, ok := obj.(*Cons)
	if !ok {
		return nil, fmt.Errorf("cons expected, but got %v", obj)
	}
	return c.car, nil
}

func Cdr(obj Object) (Object, error) {
	c, ok := obj.(*Cons)
	if !ok {
		return nil, fmt.Errorf("cons expected, but got %v", obj)
	}
	return c.cdr, nil
}

func FromBool(b bool) Object {
	if b {
		return true
	} else {
		return nil
	}
}

func NewCons(car Object, cdr Object) Object {
	return &Cons{car, cdr}
}

func NewFunc(code Code, env *Env) *Func {
	return &Func{code, env}
}

func ListToSlice(obj Object) ([]Object, Object, error) {
	if obj == nil {
		return nil, nil, nil
	}
	if IsAtom(obj) {
		return nil, nil, fmt.Errorf("cons expected, but got %v", obj)
	}
	var ret []Object
	c := obj.(*Cons)
	for {
		ret = append(ret, c.car)
		obj := c.cdr
		switch o := obj.(type) {
		case nil:
			return ret, nil, nil
		case *Cons:
			c = o
		default:
			return ret, o, nil
		}
	}
}

func listToString(obj Object) string {
	elems, improper, _ := ListToSlice(obj)
	var sb strings.Builder
	sb.WriteRune('(')
	for i, elem := range elems {
		s := ToString(elem)
		sb.WriteString(s)
		if i < len(elems)-1 {
			sb.WriteRune(' ')
		}
	}
	if improper != nil {
		sb.WriteString(" . ")
		sb.WriteString(ToString(improper))
	}
	sb.WriteRune(')')
	return sb.String()
}

func ToString(obj Object) string {
	switch obj := obj.(type) {
	case nil:
		return "nil"
	case bool:
		if !obj {
			panic(fmt.Sprintf("unknown type of object found: %v", obj))
		}
		return "t"
	case int:
		return strconv.Itoa(obj)
	case string:
		return fmt.Sprintf("\"%s\"", obj)
	case *Symbol:
		return obj.name
	case *Cons:
		return listToString(obj)
	case *Func:
		return "#<func>"
	default:
		panic(fmt.Sprintf("unknown type of object found: %v", obj))
	}
}
