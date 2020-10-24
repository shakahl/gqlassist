package utils

import (
	"bufio"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"
)

func StringFirstLetter(s string) string {
	return s[0:1]
}

func StringUpperCaseFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func StringLowerCaseFirst(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func IsFirstIndexOf(index int, list interface{}) {

}

func RemoveCStyleComments(content []byte) []byte {
	// http://blog.ostermiller.org/find-comment
	ccmt := regexp.MustCompile(`/\*([^*]|[\r\n]|(\*+([^*/]|[\r\n])))*\*+/`)
	return ccmt.ReplaceAll(content, []byte(""))
}

func RemoveCppStyleComments(content []byte) []byte {
	cppcmt := regexp.MustCompile(`//.*`)
	return cppcmt.ReplaceAll(content, []byte(""))
}

func RemoveCAndCppComments(srcpath, dstpath string) {
	b, err := ioutil.ReadFile(srcpath)
	if err != nil {
		panic(err)
	}
	b = RemoveCppStyleComments(RemoveCStyleComments(b))
	err = ioutil.WriteFile(dstpath, b, 0644)
	if err != nil {
		panic(err)
	}
}

func StripComment(source string, commentChars string, charPairs string) (string, error) {
	const (
		defaultCommentChars     = "#;"
		defaultCommentCharPairs = "/"
	)

	if commentChars == "" && charPairs == "" {
		commentChars = defaultCommentChars
		charPairs = defaultCommentCharPairs
	}

	sreader := strings.NewReader(source)
	scanner := bufio.NewScanner(sreader)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if pos := strings.IndexAny(line, commentChars+charPairs); pos >= 0 {
			tc := line[pos : pos+1]
			nd := strings.IndexAny(commentChars, tc)
			if nd == -1 || (len(line) >= pos+2 && tc == line[pos+1:pos+2]) {
				line = strings.TrimRightFunc(line[:pos], unicode.IsSpace)
			}
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	out := strings.Join(lines, "\n")

	return out, nil
}
