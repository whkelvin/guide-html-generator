package main

import (
	"fmt"
	. "html_generator/pkg/generator"
	"html_generator/pkg/parser"
	"os"
)

func main() {
	series := parser.ParseFolder("./assets/input")

	bookNames := []string{}
	for i := 0; i < len(series.Books); i++ {
		bookNames = append(bookNames, series.Books[i].Name)
		chapters := []Chapter{}
		for j := 0; j < len(series.Books[i].Chapters); j++ {
			chapters = append(chapters, Chapter{
				Filename: series.Books[i].Chapters[j].Filename,
				Name:     series.Books[i].Chapters[j].Name,
			})
			headings := []Heading{}
			for k := 0; k < len(series.Books[i].Chapters[j].Headings); k++ {
				headings = append(headings, Heading{
					Name: series.Books[i].Chapters[j].Headings[k].Name,
					Type: series.Books[i].Chapters[j].Headings[k].Type,
					Url:  series.Books[i].Chapters[j].Headings[k].Url,
				})
			}
			contentTmplInput := ContentTemplateInput{
				BasePath:      "./assets/mp3/",
				Headings:      headings,
				Next:          series.Books[i].Chapters[j].Next,
				Prev:          series.Books[i].Chapters[j].Prev,
				Title:         series.Books[i].Chapters[j].Name,
				RangeAudioUrl: series.Books[i].Chapters[j].RangeAudioUrl,
			}
			contentTmplOut := GenerateContent(contentTmplInput)
			filename := fmt.Sprintf("冊%v表%v.html", i+1, j+1)
			saveToFile(filename, contentTmplOut)
		}
		bookTemplateInput := BookTOCTemplateInput{
			Chapters: chapters,
		}
		bookTmplOut := GenerateBookTOC(bookTemplateInput)
		filename := fmt.Sprintf("冊%v.html", i+1)
		saveToFile(filename, bookTmplOut)
	}
	booksTemplateInput := BooksTOCTemplateInput{
		BookNames: bookNames}
	booksTmplOut := GenerateBooksTOC(booksTemplateInput)
	saveToFile("index.html", booksTmplOut)
}

func saveToFile(filename string, content string) {
	name := "./out/" + filename
	file, err := os.Create(name)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	file.WriteString(content)
}
