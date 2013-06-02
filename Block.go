package cssminify

import (
	"bufio"
	"os"
	"sync"
)

type Block struct {
	selector []byte
	pairs    []Pair
}

func Blocks(cb chan Block, file string, wg sync.WaitGroup) {
	cf := make(chan byte)

	wg.Add(1)
	go readFile(cf, file, wg)
	wg.Add(1)
	go parse(cf, cb, wg)
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

func readFile(cf chan byte, root string, wg sync.WaitGroup) {
	file, err := os.Open(root)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for b, err := reader.ReadByte(); err != nil; b, err = reader.ReadByte() {
		cf <- b
	}
	wg.Done()
}

func parse(cf chan byte, cb chan Block, wg sync.WaitGroup) {
	var letter byte
	state := new(State)
	for letter = <-cf; letter != 0; letter = <-cf {
		state.parse(cf, cb)
	}
	wg.Done()
}
