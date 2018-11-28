package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const gemojiDBJsonURL = "https://raw.githubusercontent.com/CodeFreezr/emojo/master/db/v5/emoji-v5.json"

type GemojiEmoji struct {
	No          int    `json:"No"`
	Emoji       string `json:"Emoji"`
	Category    string `json:"Category"`
	SubCategory string `json:"SubCategory"`
	Unicode     string `json:"Unicode"`
	Name        string `json:"Name"`
	Tags        string `json:"Tags"`
	Shortcode   string `json:"Shortcode"`
}

type TemplateData struct {
	PkgName string
	CodeMap map[string]string
}

const templateMapCode = `
package {{.PkgName}}

// NOTE: THIS FILE WAS PRODUCED BY THE
// EMOJICODEMAP CODE GENERATION TOOL (github.com/kyokomi/generateEmojiCodeMap)
// DO NOT EDIT

// Mapping from character to concrete escape code.
var emojiCodeMap = map[string]string{
	{{range $key, $val := .CodeMap}}":{{$key}}:": {{$val}},
{{end}}
}
`

var pkgName string
var fileName string

func init() {
	log.SetFlags(log.Llongfile)

	flag.StringVar(&pkgName, "pkg", "main", "output package")
	flag.StringVar(&fileName, "o", "emoji_codemap.go", "output file")
	flag.Parse()
}

func main() {

	codeMap, err := generateJson(pkgName)
	if err != nil {
		log.Fatalln(err)
	}

	os.Remove(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if _, err := file.Write(codeMap); err != nil {
		log.Fatalln(err)
	}
}

func generateJson(pkgName string) ([]byte, error) {

	// Read Emoji file

	res, err := http.Get(gemojiDBJsonURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	emojiFile, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var gs []GemojiEmoji
	if err := json.Unmarshal(emojiFile, &gs); err != nil {
		return nil, err
	}

	emojiCodeMap := make(map[string]string)
	for _, gemoji := range gs {
		emojiCodeMap[strings.Replace(gemoji.Shortcode, ":", "", 2)] = fmt.Sprintf("%+q", gemoji.Emoji)
	}

	// Template GenerateSource

	var buf bytes.Buffer
	t := template.Must(template.New("template").Parse(templateMapCode))
	if err := t.Execute(&buf, TemplateData{pkgName, emojiCodeMap}); err != nil {
		return nil, err
	}

	// gofmt

	bts, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(string(buf.Bytes()))
		return nil, fmt.Errorf("gofmt: %s", err)
	}

	return bts, nil
}
