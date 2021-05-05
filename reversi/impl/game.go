package reversi

import "math/rand"

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

type Strategy func(*Board, Cell) *Pos
type Game struct {
	board      *Board
	turn       Cell
	skips      map[Cell]int
	strategies map[Cell]Strategy
}

func NewGame(board *Board, turn Cell, strategies map[Cell]Strategy) *Game {
	return &Game{
		board:      board,
		turn:       turn,
		skips:      map[Cell]int{},
		strategies: strategies,
	}
}

func (b *Board) collectFlippables(pos *Pos, cell Cell) [][]Pos {
	ret := make([][]Pos, 0)
	for _, dir := range dirs {
		p := Pos{pos.Y, pos.X}
		flippables := make([]Pos, 0)
		for {
			p.Y += dir.dy
			p.X += dir.dx
			if c, ok := b.GetCell(&p); !ok {
				break
			} else if c == OpponentOf(cell) {
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

func (b *Board) IsAvailable(pos *Pos, cell Cell) bool {
	if !b.IsValidPos(pos) {
		return false
	}
	if c := b.MustGetCell(pos); c != Empty {
		return false
	}
	if len(b.collectFlippables(pos, cell)) == 0 {
		return false
	}
	return true
}

func (b *Board) collectAvailablePositions(cell Cell) []Pos {
	ret := make([]Pos, 0)
	b.ForEachPos(func(pos *Pos) {
		if b.IsAvailable(pos, cell) {
			ret = append(ret, *pos)
		}
	})
	return ret
}

func (game *Game) BoardContent() string {
	return game.board.String()
}

func (game *Game) Put(pos *Pos, cell Cell) {
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
	for _, chunk := range game.board.collectFlippables(pos, turn) {
		for _, pos := range chunk {
			game.board.MustSetCell(&pos, turn)
		}
	}
	game.turn = OpponentOf(turn)
	return false
}

func (game *Game) isPlayable(cell Cell) bool {
	return len(game.board.collectAvailablePositions(cell)) > 0
}

func (game *Game) IsOver() bool {
	return game.board.IsFull() || !game.isPlayable(game.turn)
}

func (game *Game) Scores() map[Cell]int {
	ret := map[Cell]int{}
	game.board.ForEachPos(func(pos *Pos) {
		if c := game.board.MustGetCell(pos); c != Empty {
			ret[c]++
		}
	})
	return ret
}

func (game *Game) Winner() Cell {
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
			return Empty //FIXME: represents draw
		}
	}
	return opponent
}

func TopLeftPossibleStrategy(b *Board, turn Cell) *Pos {
	if available := b.collectAvailablePositions(turn); len(available) > 0 {
		return &available[0]
	}
	return nil
}

func RandomPossibleStrategy(b *Board, turn Cell) *Pos {
	if available := b.collectAvailablePositions(turn); len(available) > 0 {
		return &available[rand.Intn(len(available))]
	}
	return nil
}
