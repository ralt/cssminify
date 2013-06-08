package cssminify

import (
	"fmt"
	"regexp"
	"strings"
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
	selectors := strings.Split(selector, ",")
	for i, sel := range selectors {
		fmt.Printf("%s", minifySelector(sel))
		if i != len(selectors)-1 {
			fmt.Print(",")
		}
	}
}

func minifySelector(sel string) string {
	return cleanSpaces(sel)
}

func showPropVals(pairs []Pair) {
	for i, pair := range pairs {
		fmt.Printf("%s:%s", minifyProp(string(pair.property)), minifyVal(string(pair.value)))

		// Let's gain some space: semicolons are optional for the last value
		if i != len(pairs)-1 {
			fmt.Print(";")
		}
	}
}

func minifyProp(property string) string {
	return cleanSpaces(property)
}

func cleanSpaces(str string) string {
	str = strings.TrimSpace(str)
	re := regexp.MustCompile(`\s\s`)
	for str = re.ReplaceAllString(str, " "); re.Find([]byte(str)) != nil; {
		str = re.ReplaceAllString(str, " ")
	}
	return str
}
