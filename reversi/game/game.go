package game

import (
	"math/rand"

	"github.com/athos/go-playground/reversi/board"
)

type Turn int

const (
	Neither Turn = 0
	White        = Turn(board.White)
	Black        = Turn(board.Black)
)

func OpponentOf(turn Turn) Turn {
	return Turn(board.OpponentOf(board.Cell(turn)))
}

type Strategy func(*board.Board, board.Cell) *board.Pos
type Game struct {
	board      *board.Board
	turn       Turn
	skips      map[Turn]int
	strategies map[Turn]Strategy
}

func NewGame(b *board.Board, turn Turn, strategies map[Turn]Strategy) *Game {
	return &Game{
		board:      b,
		turn:       turn,
		skips:      map[Turn]int{},
		strategies: strategies,
	}
}

func (game *Game) Board() *board.Board {
	return game.board
}

func collectAvailablePositions(b *board.Board, turn Turn) []board.Pos {
	ret := make([]board.Pos, 0)
	b.ForEachPos(func(pos *board.Pos) {
		if b.IsAvailable(pos, board.Cell(turn)) {
			ret = append(ret, *pos)
		}
	})
	return ret
}

func (game *Game) Step() bool {
	turn := game.turn
	strategy := game.strategies[turn]
	cell := board.Cell(turn)
	pos := strategy(game.board, cell)
	if pos == nil {
		return true
	}
	game.board.MustPut(pos, cell)
	game.turn = OpponentOf(turn)
	return false
}

func (game *Game) isPlayable(turn Turn) bool {
	return len(collectAvailablePositions(game.board, turn)) > 0
}

func (game *Game) IsOver() bool {
	return game.board.IsFull() || !game.isPlayable(game.turn)
}

func (game *Game) Scores() map[Turn]int {
	ret := map[Turn]int{}
	game.board.ForEachPos(func(pos *board.Pos) {
		if c := game.board.MustGetCell(pos); c != board.Empty {
			ret[Turn(c)]++
		}
	})
	return ret
}

func (game *Game) Winner() Turn {
	turn := game.turn
	opponent := OpponentOf(turn)
	if game.board.IsFull() {
		scores := game.Scores()
		switch {
		case scores[turn] > scores[opponent]:
			return turn
		case scores[turn] < scores[opponent]:
			return opponent
		default:
			return Neither
		}
	}
	return opponent
}

func TopLeftPossibleStrategy(b *board.Board, turn board.Cell) *board.Pos {
	if available := collectAvailablePositions(b, Turn(turn)); len(available) > 0 {
		return &available[0]
	}
	return nil
}

func RandomPossibleStrategy(b *board.Board, turn board.Cell) *board.Pos {
	if available := collectAvailablePositions(b, Turn(turn)); len(available) > 0 {
		return &available[rand.Intn(len(available))]
	}
	return nil
}
