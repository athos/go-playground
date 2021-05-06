package game

import (
	"math/rand"

	"github.com/athos/go-playground/reversi/board"
)

var (
	dirs = []struct {
		dy, dx int
	}{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
)

type Strategy func(*board.Board, board.Cell) *board.Pos
type Game struct {
	board      *board.Board
	turn       board.Cell
	skips      map[board.Cell]int
	strategies map[board.Cell]Strategy
}

func NewGame(b *board.Board, turn board.Cell, strategies map[board.Cell]Strategy) *Game {
	return &Game{
		board:      b,
		turn:       turn,
		skips:      map[board.Cell]int{},
		strategies: strategies,
	}
}

func collectFlippables(b *board.Board, pos *board.Pos, cell board.Cell) [][]board.Pos {
	ret := make([][]board.Pos, 0)
	for _, dir := range dirs {
		p := board.Pos{Y: pos.Y, X: pos.X}
		flippables := make([]board.Pos, 0)
		for {
			p.Y += dir.dy
			p.X += dir.dx
			if c, ok := b.GetCell(&p); !ok {
				break
			} else if c == board.OpponentOf(cell) {
				flippables = append(flippables, p)
			} else if c == cell && len(flippables) > 0 {
				ret = append(ret, flippables)
			} else { // c == Empty || len(flippables) == 0
				break
			}
		}
	}
	return ret
}

func IsAvailable(b *board.Board, pos *board.Pos, cell board.Cell) bool {
	if !b.IsValidPos(pos) {
		return false
	}
	if c := b.MustGetCell(pos); c != board.Empty {
		return false
	}
	if len(collectFlippables(b, pos, cell)) == 0 {
		return false
	}
	return true
}

func collectAvailablePositions(b *board.Board, cell board.Cell) []board.Pos {
	ret := make([]board.Pos, 0)
	b.ForEachPos(func(pos *board.Pos) {
		if IsAvailable(b, pos, cell) {
			ret = append(ret, *pos)
		}
	})
	return ret
}

func (game *Game) BoardContent() string {
	return game.board.String()
}

func (game *Game) Put(pos *board.Pos, cell board.Cell) {
	game.board.MustSetCell(pos, cell)
}

func (game *Game) Step() bool {
	turn := game.turn
	strategy := game.strategies[turn]
	pos := strategy(game.board, turn)
	if pos == nil {
		return true
	}
	game.Put(pos, turn)
	for _, chunk := range collectFlippables(game.board, pos, turn) {
		for _, pos := range chunk {
			game.board.MustSetCell(&pos, turn)
		}
	}
	game.turn = board.OpponentOf(turn)
	return false
}

func (game *Game) isPlayable(cell board.Cell) bool {
	return len(collectAvailablePositions(game.board, cell)) > 0
}

func (game *Game) IsOver() bool {
	return game.board.IsFull() || !game.isPlayable(game.turn)
}

func (game *Game) Scores() map[board.Cell]int {
	ret := map[board.Cell]int{}
	game.board.ForEachPos(func(pos *board.Pos) {
		if c := game.board.MustGetCell(pos); c != board.Empty {
			ret[c]++
		}
	})
	return ret
}

func (game *Game) Winner() board.Cell {
	turn := game.turn
	opponent := board.OpponentOf(turn)
	if game.board.IsFull() {
		scores := game.Scores()
		switch {
		case scores[turn] > scores[opponent]:
			return turn
		case scores[turn] < scores[opponent]:
			return opponent
		default:
			return board.Empty //FIXME: represents draw
		}
	}
	return opponent
}

func TopLeftPossibleStrategy(b *board.Board, turn board.Cell) *board.Pos {
	if available := collectAvailablePositions(b, turn); len(available) > 0 {
		return &available[0]
	}
	return nil
}

func RandomPossibleStrategy(b *board.Board, turn board.Cell) *board.Pos {
	if available := collectAvailablePositions(b, turn); len(available) > 0 {
		return &available[rand.Intn(len(available))]
	}
	return nil
}
