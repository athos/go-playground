package game

import (
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

type Game struct {
	board      *board.Board
	turn       Turn
	skipLimit  int
	skips      map[Turn]int
	strategies map[Turn]Strategy
}

type GameResult struct {
	Winner Turn
	Scores map[Turn]int
	Skips  map[Turn]int
}

func NewGame(b *board.Board, turn Turn, strategies map[Turn]Strategy) *Game {
	return &Game{
		board:      b,
		turn:       turn,
		skipLimit:  2,
		skips:      map[Turn]int{},
		strategies: strategies,
	}
}

func (game *Game) Board() *board.Board {
	return game.board
}

func collectAvailablePositions(b *board.Board, turn Turn) (ret []board.Pos) {
	b.ForEachPos(func(pos *board.Pos) {
		if b.IsAvailable(pos, board.Cell(turn)) {
			ret = append(ret, *pos)
		}
	})
	return ret
}

func (game *Game) Step() {
	turn := game.turn
	strategy := game.strategies[turn]
	cell := board.Cell(turn)
	if pos := strategy(game.board, cell); pos == nil {
		game.skips[turn]++
	} else {
		game.board.MustPut(pos, cell)
	}
	game.turn = OpponentOf(turn)
}

func (game *Game) isPlayable(turn Turn) bool {
	return len(collectAvailablePositions(game.board, turn)) > 0
}

func (game *Game) IsOver() bool {
	if game.board.IsFull() {
		return true
	}
	if game.skips[White] >= game.skipLimit || game.skips[Black] >= game.skipLimit {
		return true
	}
	return !game.isPlayable(game.turn) && !game.isPlayable(OpponentOf(game.turn))
}

func (game *Game) scores() map[Turn]int {
	dist := game.board.Distribution()
	return map[Turn]int{
		White: dist[board.White],
		Black: dist[board.Black],
	}
}

func (game *Game) winner() Turn {
	turn := game.turn
	opponent := OpponentOf(turn)
	switch {
	case game.skips[turn] >= game.skipLimit:
		return opponent
	case game.skips[opponent] >= game.skipLimit:
		return turn
	}
	scores := game.scores()
	switch {
	case scores[turn] > scores[opponent]:
		return turn
	case scores[turn] < scores[opponent]:
		return opponent
	default:
		return Neither
	}
}

func (game *Game) Result() *GameResult {
	return &GameResult{
		Winner: game.winner(),
		Scores: game.scores(),
		Skips:  game.skips,
	}
}
