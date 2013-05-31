package cssminify

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func Files() []Cssfile {
	return findFiles()
}

func findFiles() []Cssfile {
	files := make([]Cssfile, 0)
	filepath.Walk(".",
		func(root string, info os.FileInfo, err error) error {

			matched, _ := regexp.MatchString(".css$", root)
			if matched {
				file := Cssfile{pathname: root, content: readFile(root)}
				files = append(files, file)
			}
			return err
		})
	return files
}

func readFile(root string) string {
	content, err := ioutil.ReadFile(root)
	if err != nil {
		panic(err)
	}
	return string(content)
}

type Cssfile struct {
	pathname string
	content  string
}
