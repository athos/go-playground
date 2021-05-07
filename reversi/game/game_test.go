package game

import (
	"testing"

	"github.com/athos/go-playground/reversi/board"
	"github.com/stretchr/testify/assert"
)

const (
	W = board.White
	B = board.Black
)

func pos(y, x int) *board.Pos {
	return &board.Pos{Y: y, X: x}
}

func setupBoard(rows, cols int, args ...interface{}) *board.Board {
	b := board.NewBoard(rows, cols)
	for len(args) > 0 {
		pos := args[0].(*board.Pos)
		cell := args[1].(board.Cell)
		b.MustSetCell(pos, cell)
		args = args[2:]
	}
	return b
}

func TestCollectAvailablePositions(t *testing.T) {
	b := setupBoard(4, 4,
		pos(1, 1), W,
		pos(1, 2), B,
		pos(2, 1), B,
		pos(2, 2), W,
	)
	assert.Equal(t, []board.Pos{
		*pos(0, 2),
		*pos(1, 3),
		*pos(2, 0),
		*pos(3, 1),
	}, collectAvailablePositions(b, White))
	assert.Equal(t, []board.Pos{
		*pos(0, 1),
		*pos(1, 0),
		*pos(2, 3),
		*pos(3, 2),
	}, collectAvailablePositions(b, Black))
	b.MustPut(pos(2, 0), W)
	assert.Equal(t, []board.Pos{
		*pos(0, 2),
		*pos(0, 3),
		*pos(1, 3),
	}, collectAvailablePositions(b, White))
	assert.Equal(t, []board.Pos{
		*pos(1, 0),
		*pos(3, 0),
		*pos(3, 2),
	}, collectAvailablePositions(b, Black))
}

func TestStep(t *testing.T) {
	b := setupBoard(4, 4,
		pos(1, 1), W,
		pos(1, 2), B,
		pos(2, 1), B,
		pos(2, 2), W,
	)
	g := NewGame(b, White, map[Turn]Strategy{
		White: func(_ *board.Board, _ board.Cell) *board.Pos { return pos(2, 0) },
		Black: func(_ *board.Board, _ board.Cell) *board.Pos { return pos(3, 2) },
	})
	g.Step()
	assert.Equal(t, `+-+-+-+-+
| | | | |
| |o|x| |
|o|o|o| |
| | | | |
+-+-+-+-+`, g.Board().String())
	g.Step()
	assert.Equal(t, `+-+-+-+-+
| | | | |
| |o|x| |
|o|o|x| |
| | |x| |
+-+-+-+-+`, g.Board().String())
}

func TestIsPlayable(t *testing.T) {
	b1 := setupBoard(4, 4,
		pos(1, 1), W,
		pos(1, 2), B,
		pos(2, 1), B,
		pos(2, 2), W,
	)
	g1 := NewGame(b1, White, map[Turn]Strategy{})
	assert.Equal(t, true, g1.isPlayable(White))
	assert.Equal(t, true, g1.isPlayable(Black))

	b2 := setupBoard(3, 3,
		pos(1, 0), W,
		pos(1, 1), W,
		pos(2, 0), B,
		pos(2, 1), W,
	)
	g2 := NewGame(b2, White, map[Turn]Strategy{})
	assert.Equal(t, false, g2.isPlayable(White))
	assert.Equal(t, true, g2.isPlayable(Black))
}

func TestIsOver(t *testing.T) {
	b1 := setupBoard(3, 3,
		pos(1, 0), B,
		pos(1, 1), B,
		pos(2, 0), W,
		pos(2, 1), B,
	)
	g1 := NewGame(b1, White, map[Turn]Strategy{})
	assert.Equal(t, false, g1.IsOver())
	b2 := setupBoard(2, 2,
		pos(0, 0), W,
		pos(0, 1), B,
		pos(1, 0), W,
		pos(1, 1), B,
	)
	g2 := NewGame(b2, White, map[Turn]Strategy{})
	assert.Equal(t, true, g2.IsOver())
	b3 := setupBoard(2, 2,
		pos(0, 0), W,
		pos(1, 0), B,
		pos(1, 1), W,
	)
	g3 := NewGame(b3, White, map[Turn]Strategy{})
	assert.Equal(t, true, g3.IsOver())
}
