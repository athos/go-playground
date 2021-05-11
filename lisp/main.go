package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	lisp "github.com/athos/go-playground/lisp/impl"
)

func step(input string) (lisp.Object, error) {
	obj, err := lisp.ReadFromString(input)
	if err != nil {
		return nil, err
	}
	code, err := lisp.Compile(obj)
	if err != nil {
		return nil, err
	}
	vm := lisp.NewVM(code)
	v, err := vm.Run()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		v, err := step(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		fmt.Println(lisp.ToString(v))
	}
}
