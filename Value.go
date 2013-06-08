package cssminify

import (
	"regexp"
	"strings"
)

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
	for _, hex := range matches {
		if isFull(hex) {
			r := strings.NewReplacer(hex, newHex(hex))
			value = r.Replace(value)
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
	re := regexp.MustCompile(`url\((.*)\)`)
	matches := re.FindStringSubmatch(value)
	if len(matches) != 2 {
		return value
	}

	img := removeQuotes(matches[1])
	bytedImg := []byte(img)
	if string(bytedImg[0:3]) == "http" {
		img = getWebImg(img)
	} else {
		img = getLocalImg(img)
	}
	return img
}

/**
 * If there are quotes or double quotes around the url, take them off.
 */
func removeQuotes(img string) string {
	bytedImg := []byte(img)
	first := bytedImg[0]
	if first == '\'' || first == '"' {
		return string(bytedImg[1 : len(bytedImg)-1])
	}
	return img
}

func getWebImg(img string) string {
	return img
}

func getLocalImg(img string) string {
	return img
}
