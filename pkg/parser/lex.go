package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// itemType identifies the type of lex items.
type itemType int

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

type item struct {
	typ itemType // The type of this item.
	pos Pos      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
}

const (
	eof        = -1
	spaceChars = " \t\r\n" // These are the space characters defined by Go itself.
)

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemText       // plain text
	itemVariable   // variable inbetween '${{' and '}}'
	itemLeftDelim  // left action delimiter '${{'
	itemRightDelim // right action delimiter '}}'
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	input   string // the string being lexed
	pos     Pos    // current position in the input
	start   Pos    // start position of this item
	width   Pos    // width of last rune read from input
	lastPos Pos    // position of most recent item returned by nextItem

	items chan item // channel of lexed items
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += Pos(w)
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	item := item{t, l.start, l.input[l.start:l.pos]}

	l.items <- item
	l.lastPos = l.start
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

// lexText scans until encountering with "$" or an opening action delimiter, "${".
func lexText(l *lexer) stateFn {
	if x := strings.Index(l.input[l.pos:], "${{"); x >= 0 {
		l.pos += Pos(x)
		if l.pos > l.start {
			l.emit(itemText)
		}

		l.pos += 3
		l.emit(itemLeftDelim)

		// get closing brace
		x2 := strings.Index(l.input[l.pos:], "}}")
		if x2 < 0 {
			return l.errorf("closing brace expected")
		}

		return lexLeftDelim(l.pos + Pos(x2))
	}

	l.pos = Pos(len(l.input))
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexLeftDelim(closingPos Pos) func(l *lexer) stateFn {
	return func(l *lexer) stateFn {
		switch r := l.next(); {
		case r == ' ' || r == '\t':
			l.ignore()
			return lexLeftDelim(closingPos)
		case unicode.IsLetter(r):
			return lexVariable(closingPos, false)
		case r == eof || isEndOfLine(r):
			return l.errorf("closing brace expected")
		default:
			return l.errorf("unexpected character: %q", r)
		}
	}
}

func lexVariable(closingPos Pos, hasSpace bool) func(l *lexer) stateFn {
	return func(l *lexer) stateFn {
		switch r := l.next(); {
		case l.pos > closingPos:
			l.pos--
			if l.pos > l.start {
				l.emit(itemVariable)
			}
			l.pos += 2
			l.emit(itemRightDelim)
			return lexText
		case r == ' ' || r == '\t':
			l.pos--
			if l.pos > l.start {
				l.emit(itemVariable)
			}
			l.pos++
			l.ignore()
			return lexVariable(closingPos, true)
		case r == eof || isEndOfLine(r):
			return l.errorf("closing brace expected")
		case hasSpace:
			return l.errorf("unexpected space character")
		default:
			return lexVariable(closingPos, hasSpace)
		}
	}
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}
