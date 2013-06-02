package cssminify

import (
	"fmt"
)

func Minify(cb chan Block) {
	for block := <-cb; block.selector != nil; block = <-cb {
		fmt.Printf("%s\n", block.selector)
	}
}
