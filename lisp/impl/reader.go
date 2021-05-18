package lisp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var delimiters = map[rune]bool{
	'(':  true,
	')':  true,
	'\'': true,
	'"':  true,
	'.':  true,
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

func (r *Reader) readNumber(negative bool) (Object, error) {
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
	c, err := r.readRune()
	if err != nil {
		return nil, err
	}
	if c == '-' {
		next, err := r.peekRune()
		if err != nil {
			if err == io.EOF {
				return Intern("-"), nil
			}
			return nil, err
		}
		if unicode.IsDigit(next) {
			return r.readNumber(true)
		}
	}
	name, err := r.readWhile(func(c rune) bool {
		return !delimiters[c] && !unicode.IsSpace(c)
	})
	if err != nil {
		return nil, err
	}
	name = string(c) + name
	switch name {
	case "t":
		return true, nil
	case "nil":
		return nil, nil
	default:
		return Intern(name), nil
	}
}

func wrapErr(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

func (r *Reader) readList() (Object, error) {
	// discards preceding '('
	r.readRune()
	var elems []Object
	var improper Object
	for {
		r.skipWhitespaces()
		c, err := r.peekRune()
		if err != nil {
			return nil, wrapErr(err)
		}
		switch c {
		case ')':
			r.readRune()
			var ret Object = improper
			for i := range elems {
				ret = NewCons(elems[len(elems)-i-1], ret)
			}
			return ret, nil
		case '.':
			r.readRune()
			improper, err = r.Read()
			if err != nil {
				return nil, err
			}
		default:
			elem, err := r.Read()
			if err != nil {
				return nil, err
			}
			if improper != nil {
				return nil, errors.New("improper lists cannot have more than one elements on the right side of dot")
			}
			elems = append(elems, elem)
		}
	}
}

func (r *Reader) Read() (Object, error) {
	err := r.skipWhitespaces()
	if err != nil {
		return nil, err
	}
	c, err := r.peekRune()
	if err != nil {
		return nil, wrapErr(err)
	}
	switch {
	case unicode.IsDigit(c):
		return r.readNumber(false)
	case c == '(':
		return r.readList()
	case c == ')':
		return nil, errors.New("unexpected )")
	case c == '\'':
		r.readRune()
		obj, err := r.Read()
		if err != nil {
			return nil, err
		}
		return &Cons{Intern("quote"), &Cons{obj, nil}}, nil
	default:
		return r.readSymbol()
	}
}

func ReadFromString(input string) (Object, error) {
	r := NewReader(strings.NewReader(input))
	return r.Read()
}
