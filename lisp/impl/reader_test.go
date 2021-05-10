package lisp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromString(t *testing.T) {
	tests := []struct {
		in  string
		out Object
	}{
		{"nil", nil},
		{"t", true},
		{"42", 42},
		{"-123", -123},
		{"foo", &Symbol{"foo"}},
		{"(1 . 2)", &Cons{1, 2}},
		{"(+ 1 2)", &Cons{&Symbol{"+"}, &Cons{1, &Cons{2, nil}}}},
		{"(1 2 3 . 4)", &Cons{1, &Cons{2, &Cons{3, 4}}}},
		{
			"(+ (* 3 3) (* 4 4))",
			&Cons{
				&Symbol{"+"},
				&Cons{
					&Cons{&Symbol{"*"}, &Cons{3, &Cons{3, nil}}},
					&Cons{&Cons{&Symbol{"*"}, &Cons{4, &Cons{4, nil}}}, nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			obj, err := ReadFromString(tt.in)
			assert.Equal(t, tt.out, obj)
			assert.Nil(t, err)
		})
	}
}
