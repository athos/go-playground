package lisp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListToSlice(t *testing.T) {
	tests := []struct {
		title string
		in    Object
		out   []Object
	}{
		{"nil", nil, nil},
		{"(1 2 3)", &Cons{1, &Cons{2, &Cons{3, nil}}}, []Object{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			xs, err := ListToSlice(tt.in)
			assert.Equal(t, tt.out, xs)
			assert.Nil(t, err)
		})
	}
}
