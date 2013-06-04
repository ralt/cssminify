package cssminify

import (
	"fmt"
)

func Minify(blocks []Block) {
	for _, block := range blocks {
		showSelectors(string(block.selector))
		fmt.Print("{")
		showPropVals(block.pairs)
		fmt.Print("}")
	}
}

func showSelectors(selector string) {
}

func showPropVals(pairs []Pair) {
}
