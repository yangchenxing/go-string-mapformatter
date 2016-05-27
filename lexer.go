package mapformatter

import (
	"bufio"
	"bytes"
	"io"
)

type tokenType int

const (
	invalidToken tokenType = iota
	errorToken
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

func (l *lexer) next() token {
	r, _, err := l.ReadRune()
	if err == io.EOF {
		return token{
			typ: eofToken,
		}
	} else if err != nil {
		return token{
			typ: errorToken,
			err: err,
		}
	}
	if r == '%' {
		return l.scanField()
	}
	if err := l.UnreadRune(); err != nil {
		return token{
			typ: errorToken,
			err: err,
		}
	}
	return l.scanLiteral()
}

func (l *lexer) scanField() token {
	r, _, err := l.ReadRune()
	if err == io.EOF {
		return token{
			typ:  literalToken,
			text: "%",
		}
	} else if err != nil {
		return token{
			typ: errorToken,
			err: err,
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
		text: "%",
	}
}

func (l *lexer) scanKey() token {
	var buf bytes.Buffer
	for {
		r, _, err := l.ReadRune()
		if err == io.EOF {
			return token{
				typ:  literalToken,
				text: "%(" + buf.String(),
			}
		} else if err != nil {
			return token{
				typ: errorToken,
				err: err,
			}
		}
		if r == '|' {
			break
		} else if _, err := buf.WriteRune(r); err != nil {
			return token{
				typ: errorToken,
				err: err,
			}
		}
	}
	return l.scanFormat(buf.String())
}

func (l *lexer) scanFormat(text string) token {
	var buf bytes.Buffer
	for {
		r, _, err := l.ReadRune()
		if err == io.EOF {
			return token{
				typ:  literalToken,
				text: "%(" + text + "|" + buf.String(),
			}
		} else if err != nil {
			return token{
				typ: errorToken,
				err: err,
			}
		}
		if r == ')' {
			break
		}
		if _, err := buf.WriteRune(r); err != nil {
			return token{
				typ: errorToken,
				err: err,
			}
		}
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
		r, _, err := l.ReadRune()
		if err == io.EOF {
			return token{
				typ:  literalToken,
				text: buf.String(),
			}
		} else if err != nil {
			return token{
				typ: errorToken,
				err: err,
			}
		}
		if r == '%' {
			if err := l.UnreadRune(); err != nil {
				return token{
					typ: errorToken,
					err: err,
				}
			}
			break
		} else if _, err := buf.WriteRune(r); err != nil {
			return token{
				typ: errorToken,
				err: err,
			}
		}
	}
	return token{
		typ:  literalToken,
		text: buf.String(),
	}
}
