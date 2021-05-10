package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	lisp "github.com/athos/go-playground/lisp/impl"
)

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
		obj, err := lisp.ReadFromString(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		fmt.Println(lisp.ToString(obj))
	}
}
