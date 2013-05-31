package cssminify

import (
	"os"
	"path/filepath"
	"regexp"
)

func Files() []string {
	files := make([]string, 0)
	filepath.Walk(".",
		func(root string, info os.FileInfo, err error) error {
			matched, _ := regexp.MatchString(".css$", root)
			if matched {
				files = append(files, root)
			}
			return err
		})
	return files
}
