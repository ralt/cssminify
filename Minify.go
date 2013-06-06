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

func minifyVal(value string) string {
	value = cleanSpaces(value)

	// Values need special care
	value = cleanHex(value)
	value = cleanUrl(value)
	return value
}

func cleanHex(value string) string {
	re := regexp.MustCompile(`#(\d{6})`)
	matches := re.FindAllString(value, -1)
	if matches != nil {
		for _, hex := range matches {
			if isFull(hex) {
				r := strings.NewReplacer(hex, newHex(hex))
				value = r.Replace(value)
			}
		}
	}
	return value
}

func isFull(hex string) bool {
	cmp := []byte(hex)[1:]
	for _, l := range cmp {
		if cmp[0] != l {
			return false
		}
	}
	return true
}

func newHex(hex string) string {
	return string([]byte(hex)[:4])
}

func cleanUrl(value string) string {
	return value
}

func cleanSpaces(str string) string {
	str = strings.TrimSpace(str)
	re := regexp.MustCompile(`\s\s`)
	for str = re.ReplaceAllString(str, " "); re.Find([]byte(str)) != nil; {
		str = re.ReplaceAllString(str, " ")
	}
	return str
}
