package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowsAndCols(t *testing.T) {
	b23 := NewBoard(2, 3)
	assert.Equal(t, 2, b23.Rows())
	assert.Equal(t, 3, b23.Cols())
	b44 := NewBoard(4, 4)
	assert.Equal(t, 4, b44.Rows())
	assert.Equal(t, 4, b44.Cols())
}

func TestGetAndSet(t *testing.T) {
	b := NewBoard(4, 4)
	assert.Equal(t, Empty, b.MustGetCell(&Pos{1, 1}))
	assert.Equal(t, Empty, b.MustGetCell(&Pos{1, 2}))
	b.MustSetCell(&Pos{1, 1}, White)
	assert.Equal(t, White, b.MustGetCell(&Pos{1, 1}))
	assert.Equal(t, Empty, b.MustGetCell(&Pos{1, 2}))
}

func TestIsFull(t *testing.T) {
	b := NewBoard(2, 2)
	b.MustSetCell(&Pos{0, 0}, White)
	b.MustSetCell(&Pos{0, 1}, Black)
	b.MustSetCell(&Pos{1, 0}, Black)
	assert.Equal(t, false, b.IsFull())
	b.MustSetCell(&Pos{1, 0}, White)
	assert.Equal(t, false, b.IsFull())
	b.MustSetCell(&Pos{1, 1}, White)
	assert.Equal(t, true, b.IsFull())
}

func TestCollectFlippables(t *testing.T) {
	b := NewBoard(4, 4)
	b.MustSetCell(&Pos{1, 1}, White)
	b.MustSetCell(&Pos{1, 2}, White)
	b.MustSetCell(&Pos{2, 0}, Black)
	b.MustSetCell(&Pos{2, 1}, White)
	b.MustSetCell(&Pos{2, 2}, White)
	b.MustSetCell(&Pos{3, 2}, Black)
	actual := b.collectFlippables(&Pos{0, 2}, Black)
	expected := [][]Pos{
		{Pos{1, 1}},
		{Pos{1, 2}, Pos{2, 2}},
	}
	assert.Equal(t, expected, actual)
}

func TestBoardToString(t *testing.T) {
	b := NewBoard(8, 8)
	b.MustSetCell(&Pos{3, 3}, White)
	b.MustSetCell(&Pos{3, 4}, Black)
	b.MustSetCell(&Pos{4, 3}, Black)
	b.MustSetCell(&Pos{4, 4}, White)
	actual := b.String()
	expected := `+-+-+-+-+-+-+-+-+
| | | | | | | | |
| | | | | | | | |
| | | | | | | | |
| | | |o|x| | | |
| | | |x|o| | | |
| | | | | | | | |
| | | | | | | | |
| | | | | | | | |
+-+-+-+-+-+-+-+-+`
	assert.Equal(t, expected, actual)
}
