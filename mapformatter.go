package mapformatter

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
	"fmt"
)

var (
	invalidTokenError = errors.New("lexer return invalid token")
	cache             = make(map[string]*mapFormatter)
)

// Format format string with a map. This function may be fail with non-nil error.
func Format(format string, m map[string]interface{}) (string, error) {
	var err error
	mf := cache[format]
	if mf == nil {
		mf, err = newMapFormatter(format)
		if err != nil {
			return "", err
		}
		cache[format] = mf
	}
	return mf.format(m), nil
}

// FormatTwo format string with two maps. If the key does not exist in first map, it will be lookup in second map.
// This function may be fail with non-nil error.
func FormatTwo(format string, m1, m2 map[string]interface{}) (string, error) {
	var err error
	mf := cache[format]
	if mf == nil {
		mf, err = newMapFormatter(format)
		if err == nil {
			return "", err
		}
		cache[format] = mf
	}
	return mf.formatTwo(m1, m2), nil
}

// MustFormat format string with a map. It returns the format parameter while formatting failed.
func MustFormat(format string, m map[string]interface{}) string {
	text, err := Format(format, m)
	if err != nil {
		return format
	}
	return text
}

// MustFormatTwo format string with two maps. It returns the format parameter while formatting failed.
func MustFormatTwo(format string, m1, m2 map[string]interface{}) string {
	text, err := FormatTwo(format, m1, m2)
	if err != nil {
		return format
	}
	return text
}

type mapFormatter struct {
	fmt  string
	keys []string
}

func newMapFormatter(format string) (*mapFormatter, error) {
	var buf bytes.Buffer
	keys := list.New()
	l := &lexer{
		Reader: bufio.NewReader(bytes.NewBufferString(format)),
	}
	var t token
	for t = l.next(); t.typ == literalToken || t.typ == verbToken; t = l.next() {
		if t.typ == literalToken {
			buf.WriteString(t.text)
		} else {
			buf.WriteRune('%')
			buf.WriteString(t.format)
			keys.PushBack(t.text)
		}
	}

	if t.typ == eofToken {
		mf := &mapFormatter{
			fmt:  buf.String(),
			keys: make([]string, keys.Len()),
		}
		for i, k := 0, keys.Front(); k != nil; i, k = i+1, k.Next() {
			mf.keys[i] = k.Value.(string)
		}
		return mf, nil
	} else if t.typ == errorToken {
		return nil, t.err
	}
	return nil, invalidTokenError
}

func (mf *mapFormatter) format(m map[string]interface{}) string {
	params := make([]interface{}, len(mf.keys))
	for i, key := range mf.keys {
		params[i] = m[key]
	}
	return fmt.Sprintf(mf.fmt, params...)
}

func (mf *mapFormatter) formatTwo(m1, m2 map[string]interface{}) string {
	params := make([]interface{}, len(mf.keys))
	for i, key := range mf.keys {
		if v, found := m1[key]; found {
			params[i] = v
		} else {
			params[i] = m2[key]
		}
	}
	return fmt.Sprintf(mf.fmt, params...)
}
