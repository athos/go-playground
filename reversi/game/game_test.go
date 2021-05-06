package game

import (
	"testing"

	"github.com/athos/go-playground/reversi/board"
	"github.com/stretchr/testify/assert"
)

func TestCollectFlippables(t *testing.T) {
	b := board.NewBoard(4, 4)
	b.MustSetCell(&board.Pos{1, 1}, board.White)
	b.MustSetCell(&board.Pos{1, 2}, board.White)
	b.MustSetCell(&board.Pos{2, 0}, board.Black)
	b.MustSetCell(&board.Pos{2, 1}, board.White)
	b.MustSetCell(&board.Pos{2, 2}, board.White)
	b.MustSetCell(&board.Pos{3, 2}, board.Black)
	actual := CollectFlippables(&board.Pos{0, 2}, board.Black)
	expected := [][]board.Pos {
		{board.Pos{1, 1}},
		{board.Pos{1, 2}, board.Pos{2, 2}},
	}
	assert.Equal(t, expected, actual)
}
