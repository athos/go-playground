package lisp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListToSlice(t *testing.T) {
	tests := []struct {
		title    string
		in       Object
		out      []Object
		improper Object
	}{
		{"nil", nil, nil, nil},
		{"(1 . 2)", &Cons{1, 2}, []Object{1}, 2},
		{"(1 2 3)", &Cons{1, &Cons{2, &Cons{3, nil}}}, []Object{1, 2, 3}, nil},
		{
			"(1 (2) 3)",
			&Cons{1, &Cons{&Cons{2, nil}, &Cons{3, nil}}},
			[]Object{1, &Cons{2, nil}, 3},
			nil,
		},
		{
			"(1 2 3 . 4)",
			&Cons{1, &Cons{2, &Cons{3, 4}}},
			[]Object{1, 2, 3},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			xs, improper, err := ListToSlice(tt.in)
			assert.Equal(t, tt.out, xs)
			assert.Equal(t, tt.improper, improper)
			assert.Nil(t, err)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		in  Object
		out string
	}{
		{nil, "nil"},
		{true, "t"},
		{42, "42"},
		{Intern("foo"), "foo"},
		{&Cons{Intern("+"), &Cons{1, &Cons{2, nil}}}, "(+ 1 2)"},
		{
			&Cons{
				Intern("+"),
				&Cons{
					&Cons{Intern("*"), &Cons{3, &Cons{3, nil}}},
					&Cons{&Cons{Intern("*"), &Cons{4, &Cons{4, nil}}}, nil},
				},
			},
			"(+ (* 3 3) (* 4 4))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.out, func(t *testing.T) {
			assert.Equal(t, tt.out, ToString(tt.in))
		})
	}
}
