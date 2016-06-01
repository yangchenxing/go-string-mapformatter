package mapformatter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type tokenType int

const (
	invalidToken tokenType = iota
	eofToken
	literalToken
	verbToken
)

type token struct {
	typ    tokenType
	text   string
	format string
	err    error
}

type lexer struct {
	*bufio.Reader
}

func (l *lexer) readRune() (r rune, err error) {
	r, _, err = l.ReadRune()
	if err != nil && err != io.EOF {
		panic(err)
	}
	return
}

func (l *lexer) next() token {
	r, err := l.readRune()
	if err == io.EOF {
		return token{
			typ: eofToken,
		}
	}
	if r == '%' {
		return l.scanField()
	}
	l.UnreadRune()
	return l.scanLiteral()
}

func (l *lexer) scanField() token {
	r, err := l.readRune()
	if err == io.EOF {
		return token{
			typ:  literalToken,
			text: "%",
		}
	}
	if r == '%' {
		return token{
			typ:  literalToken,
			text: "%%",
		}
	} else if r == '(' {
		return l.scanKey()
	}
	return token{
		typ:  literalToken,
		text: fmt.Sprintf("%%%c", r),
	}
}

func (l *lexer) scanKey() token {
	var buf bytes.Buffer
	for {
		r, err := l.readRune()
		if err == io.EOF {
			return token{
				typ:  literalToken,
				text: "%(" + buf.String(),
			}
		}
		if r == '|' {
			break
		}
		buf.WriteRune(r)
	}
	return l.scanFormat(buf.String())
}

func (l *lexer) scanFormat(text string) token {
	var buf bytes.Buffer
	for {
		r, err := l.readRune()
		if err == io.EOF {
			return token{
				typ:  literalToken,
				text: "%(" + text + "|" + buf.String(),
			}
		}
		if r == ')' {
			break
		}
		buf.WriteRune(r)
	}
	return token{
		typ:    verbToken,
		text:   text,
		format: buf.String(),
	}
}

func (l *lexer) scanLiteral() token {
	var buf bytes.Buffer
	for {
		r, err := l.readRune()
		if err == io.EOF {
			return token{
				typ:  literalToken,
				text: buf.String(),
			}
		}
		if r == '%' {
			l.UnreadRune()
			break
		}
		buf.WriteRune(r)
	}
	return token{
		typ:  literalToken,
		text: buf.String(),
	}
}
