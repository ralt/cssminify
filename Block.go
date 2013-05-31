package cssminify

func Blocks(file Cssfile) []Block {
	return make([]Block, 1)
}

type Block struct {
	selector string
	pair     Pair
}
