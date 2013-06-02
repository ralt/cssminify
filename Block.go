package cssminify

import (
	"bufio"
	"fmt"
	"os"
)

type Block struct {
	selector []byte
	pairs    []Pair
}

func Blocks(cb chan Block, file string) {
	cf := make(chan byte)

	go readFile(cf, file)
	go parse(cf, cb)
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

func readFile(cf chan byte, root string) {
	file, err := os.Open(root)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for b, err := reader.ReadByte(); err != nil; b, err = reader.ReadByte() {
		fmt.Printf("%c\n", b)
		cf <- b
	}
}

func parse(cf chan byte, cb chan Block) {
	var letter byte
	state := new(State)
	for letter = <-cf; letter != 0; letter = <-cf {
		state.parse(cf, cb)
	}
}
