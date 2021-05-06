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

func prompt(msg string) (ret string) {
	fmt.Print(msg)
	fmt.Scanf("%s", &ret)
	return ret
}

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
	if input == "" || strings.ToLower(input) == "pass" {
		return nil, nil
	}

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
	for {
		input := prompt("Your turn. Type in your hand (eg. e6): ")
		pos, err := validateUserInput(b, cell, input)
		if err != nil {
			fmt.Println()
			showBoard(b)
			fmt.Printf("[ERROR] %s\n", err.Error())
			continue
		}
		if pos == nil {
			fmt.Println("Your turn has been passed.")
		}
		return pos
	}
}

func wrapCPUStrategy(strategy game.Strategy) game.Strategy {
	return func(b *board.Board, c board.Cell) *board.Pos {
		pos := strategy(b, c)
		if pos == nil {
			fmt.Println("CPU's turn has been passed.")
			return nil
		}
		fmt.Printf("CPU's turn: %s\n", pos.String())
		return pos
	}
}

func initGame(player game.Turn) *game.Game {
	b := board.NewBoard(8, 8)
	for _, c := range []*struct {
		pos  board.Pos
		cell board.Cell
	}{
		{board.Pos{Y: 3, X: 3}, board.Black},
		{board.Pos{Y: 3, X: 4}, board.White},
		{board.Pos{Y: 4, X: 3}, board.White},
		{board.Pos{Y: 4, X: 4}, board.Black},
	} {
		b.MustSetCell(&c.pos, c.cell)
	}
	opponent := game.OpponentOf(player)
	strategies := map[game.Turn]game.Strategy{
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

func showGameResult(g *game.Game, player game.Turn) {
	opponent := game.OpponentOf(player)
	result := g.Result()
	fmt.Printf("You: %2d (%d passes)\n", result.Scores[player], result.Skips[player])
	fmt.Printf("CPU: %2d (%d passes)\n", result.Scores[opponent], result.Skips[opponent])
	winner := result.Winner
	switch {
	case winner == player:
		fmt.Print("You win.")
	case winner == opponent:
		fmt.Print("You lose.")
	default:
		fmt.Print("Game is a draw.")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	player := game.White
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
