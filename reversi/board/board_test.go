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
	assert.Equal(t, Empty, b.MustGetCell(&Pos{Y: 1, X: 1}))
	assert.Equal(t, Empty, b.MustGetCell(&Pos{Y: 1, X: 2}))
	b.MustSetCell(&Pos{Y: 1, X: 1}, White)
	assert.Equal(t, White, b.MustGetCell(&Pos{Y: 1, X: 1}))
	assert.Equal(t, Empty, b.MustGetCell(&Pos{Y: 1, X: 2}))
}

func TestIsFull(t *testing.T) {
	b := NewBoard(2, 2)
	b.MustSetCell(&Pos{Y: 0, X: 0}, White)
	b.MustSetCell(&Pos{Y: 0, X: 1}, Black)
	b.MustSetCell(&Pos{Y: 1, X: 0}, Black)
	assert.Equal(t, false, b.IsFull())
	b.MustSetCell(&Pos{Y: 1, X: 0}, White)
	assert.Equal(t, false, b.IsFull())
	b.MustSetCell(&Pos{Y: 1, X: 1}, White)
	assert.Equal(t, true, b.IsFull())
}

func TestCollectFlippables(t *testing.T) {
	b := NewBoard(4, 4)
	b.MustSetCell(&Pos{Y: 1, X: 1}, White)
	b.MustSetCell(&Pos{Y: 1, X: 2}, White)
	b.MustSetCell(&Pos{Y: 2, X: 0}, Black)
	b.MustSetCell(&Pos{Y: 2, X: 1}, White)
	b.MustSetCell(&Pos{Y: 2, X: 2}, White)
	b.MustSetCell(&Pos{Y: 3, X: 2}, Black)
	actual := b.collectFlippables(&Pos{Y: 0, X: 2}, Black)
	expected := [][]Pos{
		{Pos{Y: 1, X: 1}},
		{Pos{Y: 1, X: 2}, Pos{Y: 2, X: 2}},
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
