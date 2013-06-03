package cssminify

import (
	"fmt"
)

func Minify(blocks []Block) {
	for _, block := range blocks {
		fmt.Printf("%s\n", block.selector)
	}
}
