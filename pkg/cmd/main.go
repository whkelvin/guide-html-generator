package main

import (
	"fmt"
	"github.com/charmbracelet/log"
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

	src = "./assets/html/wh-range-audio.js"
	dst = "./out/wh-range-audio.js"
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
				Name:     series.Books[i].Chapters[j].TOCTitle,
			})

			headings := []Heading{}
			for k := 0; k < len(series.Books[i].Chapters[j].Headings); k++ {
				headings = append(headings, Heading{
					Name:        series.Books[i].Chapters[j].Headings[k].Name,
					Type:        series.Books[i].Chapters[j].Headings[k].Type,
					Url:         series.Books[i].Chapters[j].Headings[k].Url,
					LinkBtnText: series.Books[i].Chapters[j].Headings[k].BtnName,
				})
			}

			booknames := []string{}
			for n := 0; n < series.Books[i].Chapters[j].TotalBookCount; n++ {
				booknames = append(booknames, fmt.Sprintf("冊%v目錄", n+1))
			}
			contentTmplInput := ContentTemplateInput{
				BasePath:      "./assets/mp3/",
				Headings:      headings,
				Next:          series.Books[i].Chapters[j].Next,
				Prev:          series.Books[i].Chapters[j].Prev,
				Title:         series.Books[i].Chapters[j].Title,
				RangeAudioUrl: series.Books[i].Chapters[j].RangeAudioUrl,
				BookNames:     booknames,
				BookNumber:    i + 1,
			}
			contentTmplOut := GenerateContent(contentTmplInput)
			filename := fmt.Sprintf("%v.html", series.Books[i].Chapters[j].Filename)
			saveToFile(filename, contentTmplOut)
		}
		items := []TOCItem{}
		for n := 0; n < len(series.Books[i].TOC.TOCItems); n++ {
			items = append(items, TOCItem{
				Title: series.Books[i].TOC.TOCItems[n].Title,
				Type:  series.Books[i].TOC.TOCItems[n].Type,
				Url:   series.Books[i].TOC.TOCItems[n].Url,
			})
		}
		bookTemplateInput := BookTOCTemplateInput{
			Items:      items,
			BookNumber: i + 1,
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
	log.Infof("%v written to disk", filename)
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
