package game

import (
	"math/rand"

	"github.com/athos/go-playground/reversi/board"
)

type Strategy func(*board.Board, board.Cell) *board.Pos

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
