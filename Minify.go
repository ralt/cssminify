package cssminify

import (
	"fmt"
	"sync"
)

func Minify(cb chan Block, wg sync.WaitGroup) {
	for block := <-cb; block.selector != nil; block = <-cb {
		fmt.Printf("%s\n", block.selector)
	}
	wg.Done()
}
