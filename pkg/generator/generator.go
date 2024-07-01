package generator

import (
	"bytes"
	//"github.com/charmbracelet/log"
	"io"
	"os"
	"strings"
	"text/template"
)

type SeriesTOCTemplateInput struct {
	BasePath string
	TOC      []TOCItem
}

func GenerateSeriesTOC(input SeriesTOCTemplateInput) string {
	txt, err := os.ReadFile("assets/html/series-toc.html")
	if err != nil {
		panic(err.Error())
	}

	tmpl, err := template.New("SeriesTOC").Parse(string(txt))
	if err != nil {
		panic(err.Error())
	}

	var buf bytes.Buffer
	writer := io.Writer(&buf)
	err = tmpl.Execute(writer, input)
	if err != nil {
		panic(err.Error())
	}
	return buf.String()
}

type Chapter struct {
	Filename string
	Name     string
}

type TOCItem struct {
	Title string
	Type  string
	Url   string
}

type BookTOCTemplateInput struct {
	Items    []TOCItem
	BookName string
	BasePath string
}

func GenerateBookTOC(input BookTOCTemplateInput) string {
	txt, err := os.ReadFile("assets/html/book-toc.html")
	if err != nil {
		panic(err.Error())
	}

	tmpl, err := template.New("BookTOC").Parse(string(txt))
	if err != nil {
		panic(err.Error())
	}

	var buf bytes.Buffer
	writer := io.Writer(&buf)
	err = tmpl.Execute(writer, input)
	if err != nil {
		panic(err.Error())
	}
	return buf.String()
}

type Heading struct {
	Name        string
	Type        string
	Url         string
	LinkBtnText string
}

type ContentTemplateInput struct {
	Headings      []Heading
	RangeAudioUrl string
	Next          string
	Prev          string
	Title         string
	BasePath      string
	TOCUrl        string
}

func GenerateContent(input ContentTemplateInput) string {
	for i := 0; i < len(input.Headings); i++ {
		input.Headings[i].Name = strings.ReplaceAll(input.Headings[i].Name, " ", "&nbsp;")
	}

	txt, err := os.ReadFile("assets/html/content.txt")
	if err != nil {
		panic(err.Error())
	}

	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}

	tmpl, err := template.New("Content").Funcs(funcMap).Parse(string(txt))
	if err != nil {
		panic(err.Error())
	}

	var buf bytes.Buffer
	writer := io.Writer(&buf)
	err = tmpl.Execute(writer, input)
	if err != nil {
		panic(err.Error())
	}
	return buf.String()
}
