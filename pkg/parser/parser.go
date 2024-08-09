package parser

import (
	"errors"
	"strings"

	"github.com/environment-toolkit/et/pkg/resolver"
)

// A mode value is a set of flags (or 0). They control parser behavior.
type Mode int

// Mode for parser behaviour
const (
	Quick     Mode = iota // stop parsing after first error encountered and return
	AllErrors             // report all errors
)

type Parser interface {
	Parse(text string) (string, error)
}

type parser struct {
	resolverSession resolver.Session

	// parsing state;
	lex       *lexer
	token     [3]item // three-token lookahead
	peekCount int
	nodes     []Node
}

// Parse parses the given string.
func (p *parser) Parse(text string) (string, error) {
	p.lex = lex(text)

	// clean parse state
	p.nodes = make([]Node, 0)
	p.peekCount = 0
	if err := p.parse(); err != nil {
		return "", err
	}

	// resolve all values
	if err := p.resolverSession.Resolve(); err != nil {
		return "", err
	}

	var b strings.Builder
	for _, node := range p.nodes {
		s, err := node.String()
		if err != nil {
			return "", err
		}
		b.WriteString(s)
	}
	return b.String(), nil
}

// parse is the top-level parser for the template.
// It runs to EOF and return an error if something isn't right.
func (p *parser) parse() error {
	for {
		switch t := p.next(); t.typ {
		case itemEOF:
			return nil
		case itemError:
			return errors.New(t.val)
		case itemLeftDelim:
			continue
		case itemRightDelim:
			continue
		case itemVariable:
			val, err := p.resolverSession.NewValue(t.val)
			if err != nil {
				return errors.New("failed to resolve variable")
			}
			varNode := NewVariable(val)
			p.nodes = append(p.nodes, varNode)
			continue
		default:
			textNode := NewText(t.val)
			p.nodes = append(p.nodes, textNode)
		}
	}
}

// next returns the next token.
func (p *parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.nextItem()
	}
	return p.token[p.peekCount]
}

func New(resolverSession resolver.Session) Parser {
	return &parser{
		resolverSession: resolverSession,
	}
}
