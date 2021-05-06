package board

import (
	"fmt"
	"strings"
)

type Cell int

const (
	Empty Cell = iota
	White
	Black
)

func OpponentOf(c Cell) Cell {
	switch c {
	case White:
		return Black
	case Black:
		return White
	default:
		panic("Empty cell does not have opponent")
	}
}

type Pos struct{ Y, X int }

func (p *Pos) String() string {
	return fmt.Sprintf("%c%d", rune(p.X+'a'), p.Y+1)
}

type Board struct {
	rows, cols int
	remaining  int
	cells      [][]Cell
}

func NewBoard(rows, cols int) *Board {
	cells := make([][]Cell, rows)
	for i := 0; i < cols; i++ {
		cells[i] = make([]Cell, cols)
	}
	return &Board{
		rows:      rows,
		cols:      cols,
		remaining: rows * cols,
		cells:     cells,
	}
}

func (b *Board) IsFull() bool {
	return b.remaining == 0
}

func (b *Board) IsValidPos(p *Pos) bool {
	return 0 <= p.Y && p.Y < b.rows && 0 <= p.X && p.X < b.cols
}

func (b *Board) MustGetCell(p *Pos) Cell {
	return b.cells[p.Y][p.X]
}

func (b *Board) GetCell(p *Pos) (c Cell, ok bool) {
	if b.IsValidPos(p) {
		return b.MustGetCell(p), true
	}
	return
}

func (b *Board) MustSetCell(p *Pos, c Cell) {
	old := b.cells[p.Y][p.X]
	b.cells[p.Y][p.X] = c
	if old == Empty {
		b.remaining--
	}
}

func (b *Board) ForEachPos(f func(*Pos)) {
	pos := Pos{}
	for ; pos.Y < b.rows; pos.Y++ {
		for pos.X = 0; pos.X < b.cols; pos.X++ {
			f(&pos)
		}
	}
}

func (b *Board) collectFlippables(pos *Pos, cell Cell) [][]Pos {
	ret := make([][]Pos, 0)
	for _, dir := range []struct{ dy, dx int }{
		{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1},
	} {
		p := Pos{Y: pos.Y, X: pos.X}
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

func (b *Board) MustPut(pos *Pos, cell Cell) {
	b.MustSetCell(pos, cell)
	for _, chunk := range b.collectFlippables(pos, cell) {
		for _, p := range chunk {
			b.MustSetCell(&p, cell)
		}
	}
}

func (b *Board) String() string {
	sb := new(strings.Builder)
	sb.WriteString("   ")
	for i := 0; i < b.cols; i++ {
		sb.WriteRune(rune('a' + i))
		sb.WriteRune(' ')
	}
	sb.WriteString("\n  +")
	for i := 0; i < b.cols; i++ {
		sb.WriteString("-+")
	}
	pos := Pos{0, 0}
	for ; pos.Y < b.rows; pos.Y++ {
		sb.WriteString(fmt.Sprintf("\n%d |", pos.Y+1))
		for pos.X = 0; pos.X < b.cols; pos.X++ {
			switch b.MustGetCell(&pos) {
			case White:
				sb.WriteRune('o')
			case Black:
				sb.WriteRune('x')
			default:
				sb.WriteRune(' ')
			}
			sb.WriteRune('|')
		}
	}
	sb.WriteString("\n  +")
	for i := 0; i < b.cols; i++ {
		sb.WriteString("-+")
	}
	return sb.String()
}
