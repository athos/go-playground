package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/athos/go-playground/reversi/board"
	"github.com/athos/go-playground/reversi/game"
)

var (
	validInputRE = regexp.MustCompile("^([a-h])([1-8])$")
)

func showBoard(b *board.Board) {
	lines := strings.Split(b.String(), "\n")
	fmt.Print("  ")
	for i := 0; i < b.Cols(); i++ {
		fmt.Printf(" %c", rune('a'+i))
	}
	fmt.Printf("\n  %s", lines[0])
	for i, line := range lines[1 : len(lines)-1] {
		fmt.Printf("\n%d %s", i+1, line)
	}
	fmt.Printf("\n  %s\n", lines[len(lines)-1])
}

func validateUserInput(b *board.Board, cell board.Cell, input string) (*board.Pos, error) {
	match := validInputRE.FindAllStringSubmatch(input, -1)
	if len(match) == 0 {
		return nil, fmt.Errorf("invalid format of hand: %s", input)
	}

	row, col := int(match[0][2][0]-'1'), int(match[0][1][0]-'a')
	pos := &board.Pos{Y: row, X: col}
	if !b.IsAvailable(pos, cell) {
		return nil, fmt.Errorf("invalid position: %s", input)
	}
	return pos, nil
}

func userInputStrategy(b *board.Board, cell board.Cell) *board.Pos {
	var input string
	for {
		fmt.Print("Your turn. Type in your hand (eg. e6): ")
		fmt.Scanf("%s", &input)

		pos, err := validateUserInput(b, cell, input)
		if err != nil {
			fmt.Println()
			showBoard(b)
			fmt.Printf("[ERROR] %s\n", err.Error())
			continue
		}
		return pos
	}
}

func wrapCPUStrategy(strategy game.Strategy) game.Strategy {
	return func(b *board.Board, c board.Cell) *board.Pos {
		pos := strategy(b, c)
		if pos == nil {
			fmt.Println("CPU's turn has been skipped.")
			return nil
		}
		fmt.Printf("CPU's turn: %s\n", pos.String())
		return pos
	}
}

func initGame(player board.Cell) *game.Game {
	b := board.NewBoard(8, 8)
	for _, c := range []struct {
		pos  *board.Pos
		cell board.Cell
	}{
		{&board.Pos{X: 3, Y: 3}, board.Black},
		{&board.Pos{X: 3, Y: 4}, board.White},
		{&board.Pos{X: 4, Y: 3}, board.White},
		{&board.Pos{X: 4, Y: 4}, board.Black},
	} {
		b.MustSetCell(c.pos, c.cell)
	}
	opponent := board.OpponentOf(player)
	strategies := map[board.Cell]game.Strategy{
		player:   userInputStrategy,
		opponent: wrapCPUStrategy(game.RandomPossibleStrategy),
	}
	return game.NewGame(b, player, strategies)
}

func playGame(game *game.Game) {
	for {
		fmt.Println()
		showBoard(game.Board())
		if game.IsOver() {
			return
		}
		game.Step()
	}
}

func showGameResult(game *game.Game, player board.Cell) {
	opponent := board.OpponentOf(player)
	scores := game.Scores()
	fmt.Printf("You: %d\n", scores[player])
	fmt.Printf("CPU: %d\n", scores[opponent])
	winner := game.Winner()
	switch {
	case winner == player:
		fmt.Print("You win.")
	case winner == opponent:
		fmt.Print("You lose.")
	default:
		fmt.Print("Game is a draw.")
	}
}

func prompt(msg string) (ret string) {
	fmt.Print(msg)
	fmt.Scanf("%s", &ret)
	return ret
}

func main() {
	rand.Seed(time.Now().UnixNano())
	player := board.White
	for {
		game := initGame(player)
		playGame(game)
		showGameResult(game, player)
		input := strings.ToLower(prompt(" Continue? [y/N] "))
		if input != "y" {
			return
		}
	}
}
