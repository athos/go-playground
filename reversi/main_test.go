package main

import (
	"errors"
	"testing"

	"github.com/athos/go-playground/reversi/board"
	"github.com/athos/go-playground/reversi/game"
	"github.com/stretchr/testify/assert"
)

func TestValidateUserInput(t *testing.T) {
	g := initGame(game.White)
	tests := []struct {
		input string
		pos   *board.Pos
		err   error
	}{
		{"c4", &board.Pos{Y: 3, X: 2}, nil},
		{"pass", nil, nil},
		{"", nil, nil},
		{"foobar", nil, errors.New("invalid format of hand: foobar")},
		{"a1", nil, errors.New("invalid position: a1")},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			pos, err := validateUserInput(g.Board(), board.White, tt.input)
			assert.Equal(t, tt.pos, pos)
			assert.Equal(t, tt.err, err)
		})
	}
}
