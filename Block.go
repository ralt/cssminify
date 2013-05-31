package cssminify

import (
	"io/ioutil"
)

type Block struct {
	selector string
	pair     Pair
}

func Blocks(file string) []Block {
	var (
		block  Block
		blocks = make([]Block, 0)
	)

	content := readFile(file)

	block = createBlockFromFile(&content)
	for block.selector != "" {
		blocks = append(blocks, block)
		block = createBlockFromFile(&content)
	}

	return blocks
}

func createBlockFromFile(content *string) Block {
	return Block{}
}

func readFile(root string) string {
	content, err := ioutil.ReadFile(root)
	if err != nil {
		panic(err)
	}
	return string(content)
}
