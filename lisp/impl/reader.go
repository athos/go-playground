package lisp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var delimiters = map[rune]bool {
	'(': true,
	')': true,
	'\'': true,
	'"': true,
	',': true,
}

type Reader struct {
	reader *bufio.Reader
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{bufio.NewReader(reader)}
}

func (r *Reader) readRune() (rune, error) {
	c, _, err := r.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (r *Reader) unread() {
	err := r.reader.UnreadRune()
	if err != nil {
		panic(err)
	}
}

func (r *Reader) peekRune() (rune, error) {
	c, err := r.readRune()
	if err != nil {
		if err == io.EOF {
			return 0, io.ErrUnexpectedEOF
		}
		return 0, err
	}
	r.unread()
	return c, nil
}

func (r *Reader) readWhile(pred func(rune) bool) (string, error) {
	var sb strings.Builder
	for {
		c, err := r.readRune()
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
		if !pred(c) {
			r.unread()
			return sb.String(), nil
		}
		sb.WriteRune(c)
	}
}

func (r *Reader) dropWhile(pred func(rune) bool) error {
	for {
		c, err := r.readRune()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if !pred(c) {
			r.unread()
			return nil
		}
	}
}

func (r *Reader) skipWhitespaces() error {
	return r.dropWhile(unicode.IsSpace)
}

func (r *Reader) readNumber() (Object, error) {
	var negative bool
	c, err := r.peekRune()
	if err != nil {
		return nil, err
	}
	if c == '-' {
		r.readRune()
		negative = true
	}
	digits, err := r.readWhile(unicode.IsDigit)
	if err != nil {
		return nil, err
	}
	if negative {
		digits = "-" + digits
	}
	n, err := strconv.Atoi(digits)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *Reader) readSymbol() (Object, error) {
	name, err := r.readWhile(func(c rune) bool {
		return !delimiters[c] && !unicode.IsSpace(c)
	})
	if err != nil {
		return nil, err
	}
	switch name {
	case "t":
		return true, nil
	case "nil":
		return nil, nil
	default:
		return &Symbol{name}, nil
	}
}

func (r *Reader) readList() (Object, error) {
	// discards preceding '('
	r.readRune()
	var elems []Object
	for {
		c, err := r.peekRune()
		if err != nil {
			return nil, err
		}
		if c == ')' {
			r.readRune()
			var ret Object = nil
			for i := range elems {
				ret = NewCons(elems[len(elems)-i-1], ret)
			}
			return ret, nil
		}
		elem, err := r.Read()
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}
}

func (r *Reader) Read() (Object, error) {
	err := r.skipWhitespaces()
	if err != nil {
		return nil, err
	}
	c, err := r.peekRune()
	if err != nil {
		return nil, err
	}
	switch {
	case c == '-' || unicode.IsDigit(c):
		return r.readNumber()
	case c == '(':
		return r.readList()
	case c == ')':
		return nil, errors.New("unexpected )")
	default:
		return r.readSymbol()
	}
}

func ReadFromString(input string) (Object, error) {
	r := NewReader(strings.NewReader(input))
	return r.Read()
}
