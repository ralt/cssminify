package cssminify

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"regexp"
	"strings"
)

func minifyVal(value string, file string) string {
	value = cleanSpaces(value)

	// Values need special care
	value = cleanHex(value)
	value = cleanUrl(value, file)
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

func cleanUrl(value string, file string) string {
	re := regexp.MustCompile(`url\((.*)\)`)
	matches := re.FindStringSubmatch(value)
	if len(matches) != 2 {
		return value
	}

	img := removeQuotes(matches[1])
	if string([]byte(img)[0:4]) == "http" {
		img = getWebImg(img)
	} else {
		img = getLocalImg(img, file)
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
	resp, err := http.Get(img)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return writeWebUrl(body, resp)
}

func getLocalImg(img string, file string) string {
	body, err := ioutil.ReadFile(path.Dir(file) + "/" + img)
	if err != nil {
		panic(err)
	}
	return writeUrl(body, mime.TypeByExtension(img))
}

func base64Encode(data []byte) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(data)
	encoder.Close()
	return buf.String()
}

func writeWebUrl(body []byte, resp *http.Response) string {
	return writeUrl(body, resp.Header.Get("Content-Type"))
}

func writeUrl(body []byte, mime string) string {
	var buf bytes.Buffer
	buf.WriteString("url(data:")
	buf.WriteString(mime)
	buf.WriteString(";base64,")
	buf.WriteString(base64Encode(body))
	buf.WriteString(")")
	return buf.String()
}
