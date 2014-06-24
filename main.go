package main

import (
    "fmt"
	"io/ioutil"
    "log"
    "os"
    "strings"
	"unicode/utf8"
)

func main() {

	generate()
}

func generate() {
	emojiDir, err := ioutil.ReadDir("gemoji/images/emoji")
	if err != nil {
		log.Fatal(err)
	}

	currentPwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	emojiCodeMap := make(map[string]string)

	for _, emojiFileInfo := range emojiDir {
		emojiFilePath := currentPwd + "/gemoji/images/emoji/" + emojiFileInfo.Name()
		linkStr, _ := os.Readlink(emojiFilePath)
		if err != nil {
			log.Println(err)
		}

		checkCode := strings.Replace(linkStr, `unicode/`, ``, 1)
		checkCode = strings.Replace(checkCode, `.png`, ``, 1)
		count := utf8.RuneCountInString(checkCode)
		if count == 0 {
			continue
		}

		key := strings.Replace(emojiFileInfo.Name(), `.png`, ``, 1)
		code := strings.Replace(linkStr, `.png`, ``, 1)
		switch count {
		case 4: // f000
			code = strings.Replace(code, `unicode/`, `\U0000`, 1)
		case 5: // 1f000
			code = strings.Replace(code, `unicode/`, `\U000`, 1)
		case 9: // f000-f000
			code = strings.Replace(code, `unicode/`, `\U0000`, 1)
			code = strings.Replace(code, `-`, `\U0000`, 1)
		case 11: // 1f000-1f000
			code = strings.Replace(code, `unicode/`, `\U000`, 1)
			code = strings.Replace(code, `-`, `\U000`, 1)
		default:
			continue;
		}
		emojiCodeMap[key] = code
	}

	for key, value := range emojiCodeMap {
		fmt.Println(`":` + key + `:": "` + value + `",`)
	}
}
