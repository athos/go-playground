package board

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpponentOf(t *testing.T) {
	tests := []struct {
		title   string
		in, out Cell
	}{
		{"White", White, Black},
		{"Black", Black, White},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			assert.Equal(t, tt.out, OpponentOf(tt.in))
		})
	}
	assert.Panics(t, func() {
		OpponentOf(Empty)
	}, "Empty cell does not have opponent")
}

func TestPosToString(t *testing.T) {
	tests := []struct {
		pos Pos
		out string
	}{
		{Pos{Y: 0, X: 0}, "a1"},
		{Pos{Y: 3, X: 2}, "c4"},
	}
	for _, tt := range tests {
		title := fmt.Sprintf("(%d,%d) -> %s", tt.pos.Y, tt.pos.X, tt.out)
		t.Run(title, func(t *testing.T) {
			assert.Equal(t, tt.out, tt.pos.String())
		})
	}
}
func TestRowsAndCols(t *testing.T) {
	tests := []struct {
		title    string
		row, col int
	}{
		{"2x3", 2, 3},
		{"4x4", 4, 4},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			b := NewBoard(tt.row, tt.col)
			assert.Equal(t, tt.row, b.Rows())
			assert.Equal(t, tt.col, b.Cols())
		})
	}
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
