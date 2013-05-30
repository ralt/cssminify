package cssminify

func Files() []Cssfile {
	return findFiles()
}

func findFiles() []Cssfile {
	return make([]Cssfile, 1)
}

type Cssfile struct {
	pathname string
	content string
}
