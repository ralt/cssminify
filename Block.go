package cssminify

import (
	"bytes"
	"io/ioutil"
	"strings"
)

type Block struct {
	selector []byte
	pairs    []Pair
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
	STARTING_COMMENT = 0
	IN_COMMENT       = 1
	CLOSING_COMMENT  = 2
	COMMENT_CLOSED   = 3
	IN_SELECTOR      = 4
	IN_PROPERTY      = 5
	IN_VALUE         = 6
)

func Blocks(file string) []Block {
	var (
		blocks       []Block
		letter       byte
		current      []byte
		oldCurrent   []byte
		state        byte
		commentState byte
		currentBlock Block
		currentPair  Pair
	)

	content := []byte(readFile(file))

	for letter, content = stripLetter(content); letter != 0; letter, content = stripLetter(content) {
		switch letter {
		case '/':
			switch commentState {
			case CLOSING_COMMENT:
				// Since we don't keep comments
				current = oldCurrent
				commentState = COMMENT_CLOSED
			default:
				if commentState != IN_COMMENT {
					commentState = STARTING_COMMENT
					current = append(current, letter)
				}
			}
		case '*':
			switch commentState {
			case STARTING_COMMENT:
				oldCurrent = current[:len(current)-1]
				commentState = IN_COMMENT
				current = append(current, letter)
			case IN_COMMENT:
				commentState = CLOSING_COMMENT
				current = append(current, letter)
			}
		case '{':
			if commentState != IN_COMMENT {
				if state == IN_SELECTOR {
					state = IN_PROPERTY
					currentBlock.selector = current
					current = []byte{}
				} else {
					panic(NOT_IN_SELECTOR)
				}
			}
		case '}':
			if commentState != IN_COMMENT {
				if state == IN_VALUE && !bytes.Equal(nil, current) {
					state = IN_PROPERTY
					currentPair.value = current
					currentBlock.pairs = append(currentBlock.pairs, currentPair)
				}
				if state == IN_PROPERTY && strings.Trim(string(current), " ") != "" {
					current = []byte{}
					blocks = append(blocks, currentBlock)
					currentBlock = Block{}
					state = IN_SELECTOR
				} else {
					panic(NOT_AFTER_VALUE)
				}
			}
		case ':':
			if commentState != IN_COMMENT {
				if state == IN_PROPERTY && !bytes.Equal(nil, current) {
					state = IN_VALUE
					currentPair.property = current

					// Cleanup
					current = []byte{}
				} else {
					if state != IN_VALUE && state != IN_SELECTOR {
						panic(NOT_IN_PROPERTY)
					}
				}
			}
		case ';':
			if commentState != IN_COMMENT {
				if state == IN_VALUE {
					state = IN_PROPERTY
					currentPair.value = current
					currentBlock.pairs = append(currentBlock.pairs, currentPair)

					// Cleanup
					currentPair = Pair{}
					current = []byte{}
				} else {
					panic(NOT_IN_VALUE)
				}
			}
		default:
			if commentState != IN_COMMENT {
				if state == 0 {
					state = IN_SELECTOR
				}
				current = append(current, letter)
			}
		}

	}

	return blocks
}

func stripLetter(content []byte) (byte, []byte) {
	var letter byte
	if len(content) != 0 {
		letter = content[0]
		content = content[1:]
	} else {
		content = []byte{}
	}
	return letter, content
}

func readFile(root string) string {
	content, err := ioutil.ReadFile(root)
	if err != nil {
		panic(err)
	}
	return string(content)
}
