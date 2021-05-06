package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
