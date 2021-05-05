package reversi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectFlippables(t *testing.T) {
	b := NewBoard(4, 4)
	b.MustSetCell(&Pos{1, 1}, White)
	b.MustSetCell(&Pos{1, 2}, White)
	b.MustSetCell(&Pos{2, 0}, Black)
	b.MustSetCell(&Pos{2, 1}, White)
	b.MustSetCell(&Pos{2, 2}, White)
	b.MustSetCell(&Pos{3, 2}, Black)
	actual := b.CollectFlippables(&Pos{0, 2}, Black)
	expected := [][]Pos {
		{Pos{1, 1}},
		{Pos{1, 2}, Pos{2, 2}},
	}
	assert.Equal(t, expected, actual)
}
