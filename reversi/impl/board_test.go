package reversi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoardToString(t *testing.T) {
	b := NewBoard(8, 8)
	b.MustSetCell(&Pos{3, 3}, White)
	b.MustSetCell(&Pos{3, 4}, Black)
	b.MustSetCell(&Pos{4, 3}, Black)
	b.MustSetCell(&Pos{4, 4}, White)
	actual := b.String()
	expected := `+-+-+-+-+-+-+-+-+
| | | | | | | | |
| | | | | | | | |
| | | | | | | | |
| | | |o|x| | | |
| | | |x|o| | | |
| | | | | | | | |
| | | | | | | | |
| | | | | | | | |
+-+-+-+-+-+-+-+-+`
	assert.Equal(t, expected, actual)
}
