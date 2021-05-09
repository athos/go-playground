package lisp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListToSlice(t *testing.T) {
	xs, err := ListToSlice(nil)
	assert.Equal(t, []Object(nil), xs)
	assert.Equal(t, nil, err)

	ys, err := ListToSlice(&Cons{1, &Cons{2, &Cons{3, nil}}})
	assert.Equal(t, []Object{1, 2, 3}, ys)
	assert.Equal(t, nil, err)
}
