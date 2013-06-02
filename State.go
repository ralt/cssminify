package cssminify

import (
	"bytes"
	"strings"
)

type State struct {
	state        byte
	commentState byte
	current      []byte
	previous     []byte
	currentBlock Block
	currentPair  Pair
	blocks       []Block
}

// Error constants
const (
	NOT_IN_SELECTOR = "Met { while not being in selector"
	NOT_AFTER_VALUE = "Met } after a non-value"
	NOT_IN_PROPERTY = "Met : after a non-property"
	NOT_IN_VALUE    = "Met ; after a non-value"
)

// Parser constants
const (
	STARTING_COMMENT = iota
	IN_COMMENT       = iota
	CLOSING_COMMENT  = iota
	COMMENT_CLOSED   = iota
	IN_SELECTOR      = iota
	IN_PROPERTY      = iota
	IN_VALUE         = iota
)

func (s *State) parse(cf chan byte, cb chan Block) {
	letter := <-cf
	switch letter {
	case '/':
		s.slash(letter)
	case '*':
		s.star(letter)
	case '{':
		s.openBracket(letter)
	case '}':
		s.closeBracket(cb, letter)
	case ':':
		s.colon(letter)
	case ';':
		s.semicolon(letter)
	default:
		s.rest(letter)
	}
}

func (s *State) slash(letter byte) {
	switch s.commentState {
	case CLOSING_COMMENT:
		// Since we don't keep comments
		s.current = s.previous
		s.commentState = COMMENT_CLOSED
	default:
		if s.commentState != IN_COMMENT {
			s.commentState = STARTING_COMMENT
			s.current = append(s.current, letter)
		}
	}
}

func (s *State) star(letter byte) {
	switch s.commentState {
	case STARTING_COMMENT:
		s.previous = s.current[:len(s.current)-1]
		s.commentState = IN_COMMENT
		s.current = append(s.current, letter)
	case IN_COMMENT:
		s.commentState = CLOSING_COMMENT
		s.current = append(s.current, letter)
	}

}

func (s *State) openBracket(letter byte) {
	if s.commentState != IN_COMMENT {
		if s.state == IN_SELECTOR {
			s.state = IN_PROPERTY
			s.currentBlock.selector = s.current
			s.current = []byte{}
		} else {
			panic(NOT_IN_SELECTOR)
		}
	}
}

func (s *State) closeBracket(cb chan Block, letter byte) {
	if s.commentState != IN_COMMENT {
		if s.state == IN_VALUE && !bytes.Equal(nil, s.current) {
			s.state = IN_PROPERTY
			s.currentPair.value = s.current
			s.currentBlock.pairs = append(s.currentBlock.pairs, s.currentPair)
		}
		if s.state == IN_PROPERTY && strings.Trim(string(s.current), " ") != "" {
			s.current = []byte{}

			cb <- s.currentBlock

			s.currentBlock = Block{}
			s.state = IN_SELECTOR
		} else {
			panic(NOT_AFTER_VALUE)
		}
	}
}

func (s *State) colon(letter byte) {
	if s.commentState != IN_COMMENT {
		if s.state == IN_PROPERTY && !bytes.Equal(nil, s.current) {
			s.state = IN_VALUE
			s.currentPair.property = s.current

			// Cleanup
			s.current = []byte{}
		} else {
			if s.state != IN_VALUE && s.state != IN_SELECTOR {
				panic(NOT_IN_PROPERTY)
			}
		}
	}
}

func (s *State) semicolon(letter byte) {
	if s.commentState != IN_COMMENT {
		if s.state == IN_VALUE {
			s.state = IN_PROPERTY
			s.currentPair.value = s.current
			s.currentBlock.pairs = append(s.currentBlock.pairs, s.currentPair)

			// Cleanup
			s.currentPair = Pair{}
			s.current = []byte{}
		} else {
			panic(NOT_IN_VALUE)
		}
	}
}

func (s *State) rest(letter byte) {
	if s.commentState != IN_COMMENT {
		if s.state == 0 {
			s.state = IN_SELECTOR
		}
		s.current = append(s.current, letter)
	}
}
