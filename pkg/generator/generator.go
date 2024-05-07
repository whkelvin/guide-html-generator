package generator

import (
	"bytes"
	"io"
	"os"
	"strings"
	"text/template"
)

type BooksTOCTemplateInput struct {
	BookNames []string
}

func GenerateBooksTOC(input BooksTOCTemplateInput) string {
	txt, err := os.ReadFile("assets/html/books-toc.html")
	if err != nil {
		panic(err.Error())
	}

	tmpl, err := template.New("BooksTOC").Parse(string(txt))
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

type BookTOCTemplateInput struct {
	Chapters []Chapter
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
	Name string
	Type string
	Url  string
}

type ContentTemplateInput struct {
	Headings      []Heading
	RangeAudioUrl string
	Next          string
	Prev          string
	Title         string
	BasePath      string
}

func GenerateContent(input ContentTemplateInput) string {
	for i := 0; i < len(input.Headings); i++ {
		input.Headings[i].Name = strings.ReplaceAll(input.Headings[i].Name, " ", "&nbsp;")
	}

	txt, err := os.ReadFile("assets/html/content.txt")
	if err != nil {
		panic(err.Error())
	}

	tmpl, err := template.New("Content").Parse(string(txt))
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
