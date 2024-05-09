package main

import (
	"fmt"
	. "html_generator/pkg/generator"
	"html_generator/pkg/parser"
	"io"
	"os"
)

func main() {
	src := "./assets/html/wh-audio.js"
	dst := "./out/wh-audio.js"
	_, err := copyFile(src, dst)
	if err != nil {
		panic(err.Error())
	}

	src = "./assets/html/wh-raw-audio.js"
	dst = "./out/wh-raw-audio.js"
	_, err = copyFile(src, dst)
	if err != nil {
		panic(err.Error())
	}

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

func copyFile(src, dst string) (int64, error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destinationFile.Close()

	nBytes, err := io.Copy(destinationFile, sourceFile)
	if err != nil {
		return nBytes, err
	}

	return nBytes, nil
}
