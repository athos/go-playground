package lisp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Object interface{}
type Symbol struct {
	name string
}
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

func ListToSlice(obj Object) ([]Object, error) {
	if obj == nil {
		return nil, nil
	}

	c, ok := obj.(*Cons)
	if !ok {
		return nil, fmt.Errorf("cons expected, but got %v", obj)
	}
	var ret []Object
	for {
		ret = append(ret, c.car)
		obj := c.cdr
		switch o := obj.(type) {
		case nil:
			return ret, nil
		case *Cons:
			c = o
		default:
			panic("improper lists are not supported yet")
		}
	}
}

func listToString(obj Object) (string, error) {
	elems, err := ListToSlice(obj)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteRune('(')
	for i, elem := range elems {
		s, err := ToString(elem)
		if err != nil {
			return "", err
		}
		sb.WriteString(s)
		if i < len(elems) {
			sb.WriteString(", ")
		}
	}
	sb.WriteRune(')')
	return sb.String(), nil
}

func ToString(obj Object) (string, error) {
	switch obj := obj.(type) {
	case nil:
		return "nil", nil
	case int:
		return strconv.Itoa(obj), nil
	case string:
		return fmt.Sprintf("\"%s\"", obj), nil
	case *Symbol:
		return obj.name, nil
	case *Cons:
		return listToString(obj)
	case *Func:
		return "#<func>", nil
	default:
		panic(fmt.Sprintf("unknown type of object found: %v", obj))
	}
}
