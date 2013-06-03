package cssminify

import (
	"io/ioutil"
)

type Block struct {
	selector []byte
	pairs    []Pair
}

func Blocks(file string) []Block {
	var (
		letter byte
	)

	content := []byte(readFile(file))
	state := new(State)

	for letter, content = stripLetter(content); letter != 0; letter, content = stripLetter(content) {
		state.parse(letter)
	}

	return state.blocks
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
