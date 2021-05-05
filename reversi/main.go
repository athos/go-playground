package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	impl "github.com/athos/go-playground/reversi/impl"
)

var (
	validInputRE = regexp.MustCompile("^([a-h])([1-8])$")
)

func validateUserInput(b *impl.Board, cell impl.Cell, input string) (*impl.Pos, error) {
	match := validInputRE.FindAllStringSubmatch(input, -1)
	if len(match) == 0 {
		return nil, fmt.Errorf("invalid format of hand: %s", input)
	}

	row, col := int(match[0][2][0]-'1'), int(match[0][1][0]-'a')
	pos := &impl.Pos{Y: row, X: col}
	if !b.IsAvailable(pos, cell) {
		return nil, fmt.Errorf("invalid position: %s", input)
	}
	return pos, nil
}

func userInputStrategy(b *impl.Board, cell impl.Cell) *impl.Pos {
	var input string
	for {
		fmt.Print("Your turn. Type in your hand (eg. e6): ")
		fmt.Scanf("%s", &input)

		pos, err := validateUserInput(b, cell, input)
		if err != nil {
			fmt.Println()
			fmt.Println(b.String())
			fmt.Printf("[ERROR] %s\n", err.Error())
			continue
		}
		return pos
	}
}

func wrapCPUStrategy(strategy impl.Strategy) impl.Strategy {
	return func(b *impl.Board, c impl.Cell) *impl.Pos {
		pos := strategy(b, c)
		if pos == nil {
			fmt.Println("CPU's turn has been skipped.")
			return nil
		}
		fmt.Printf("CPU's turn: %s\n", pos.String())
		return pos
	}
}

func initGame(player impl.Cell) *impl.Game {
	board := impl.NewBoard(8, 8)
	opponent := impl.OpponentOf(player)
	strategies := map[impl.Cell]impl.Strategy{
		player:   userInputStrategy,
		opponent: wrapCPUStrategy(impl.RandomPossibleStrategy),
	}
	game := impl.NewGame(board, player, strategies)

	for _, c := range []struct {
		pos  *impl.Pos
		cell impl.Cell
	}{
		{&impl.Pos{X: 3, Y: 3}, impl.Black},
		{&impl.Pos{X: 3, Y: 4}, impl.White},
		{&impl.Pos{X: 4, Y: 3}, impl.White},
		{&impl.Pos{X: 4, Y: 4}, impl.Black},
	} {
		game.Put(c.pos, c.cell)
	}
	return game
}

func playGame(game *impl.Game) {
	for {
		fmt.Println()
		fmt.Println(game.BoardContent())
		if game.IsOver() {
			return
		}
		game.Step()
	}
}

func showGameResult(game *impl.Game, player impl.Cell) {
	opponent := impl.OpponentOf(player)
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
	player := impl.White
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
