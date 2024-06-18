package main

import (
	"fmt"
	. "html_generator/pkg/generator"
	"html_generator/pkg/parser"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/log"
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
		for j := 0; j < len(series.Books[i].Chapters); j++ {
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

	log.Infof("正在檢查輸出網頁")

	for i := 0; i < len(series.Books); i++ {
		log.Infof("正在檢查冊%v...", i+1)

		log.Infof("正在檢查冊%v目錄", i+1)
		for n := 0; n < len(series.Books[i].TOC.TOCItems); n++ {
			if series.Books[i].TOC.TOCItems[n].Type == "Link" {
				if doesFileExist("./out/" + series.Books[i].TOC.TOCItems[n].Url + ".html") {
					log.Infof("找到目錄標題連結: %v | %v", series.Books[i].TOC.TOCItems[n].Title, "./out/"+series.Books[i].TOC.TOCItems[n].Url+".html")
				} else {
					log.Errorf("無法找到目錄標題連結: %v | %v", series.Books[i].TOC.TOCItems[n].Title, "./out/"+series.Books[i].TOC.TOCItems[n].Url+".html")
				}
			}
		}
		for j := 0; j < len(series.Books[i].Chapters); j++ {
			log.Infof("正在檢查%v", series.Books[i].Chapters[j].Filename)

			if doesFileExist("./out/assets/mp3/本表範圍/" + series.Books[i].Chapters[j].RangeAudioUrl + ".mp3") {
				log.Info("找到本表範圍音檔!")
			} else {
				log.Errorf("無法找到本表範圍音檔: %v", "./out/assets/mp3/本表範圍/"+series.Books[i].Chapters[j].RangeAudioUrl+".mp3")
			}

			if series.Books[i].Chapters[j].Prev != "" {
				if doesFileExist("./out/" + series.Books[i].Chapters[j].Prev + ".html") {
					log.Info("找到上一表!")
				} else {
					log.Errorf("無法找到上一表: %v", "./out/"+series.Books[i].Chapters[j].Prev+".html")
				}
			}

			if series.Books[i].Chapters[j].Next != "" {
				if doesFileExist("./out/" + series.Books[i].Chapters[j].Next + ".html") {
					log.Info("找到下一表!")
				} else {
					log.Errorf("無法找到下一表: %v", "./out/"+series.Books[i].Chapters[j].Next+".html")
				}
			}

			for k := 0; k < len(series.Books[i].Chapters[j].Headings); k++ {
				if series.Books[i].Chapters[j].Headings[k].Type == "Audio" {
					log.Infof("正在檢查音檔標題：%v: %v", series.Books[i].Chapters[j].Filename, strings.TrimSpace(series.Books[i].Chapters[j].Headings[k].Name))

					if doesFileExist("./out/assets/mp3/原文/" + series.Books[i].Chapters[j].Headings[k].Url + ".mp3") {
						log.Info("找到原文音檔!")
					} else {
						log.Errorf("無法找到原文音檔: %v", "./out/assets/mp3/原文/"+series.Books[i].Chapters[j].Headings[k].Url+".mp3")
					}

					if doesFileExist("./out/assets/mp3/上下層科判/" + series.Books[i].Chapters[j].Headings[k].Url + ".mp3") {
						log.Info("找到上下層科判音檔!")
					} else {
						log.Errorf("無法找到上下層科判音檔: %v", "./out/assets/mp3/上下層科判/"+series.Books[i].Chapters[j].Headings[k].Url+".mp3")
					}

					if doesFileExist("./out/assets/mp3/各科範圍/" + series.Books[i].Chapters[j].Headings[k].Url + ".mp3") {
						log.Info("找到各科範圍音檔!")
					} else {
						log.Errorf("無法找到各科範圍音檔: %v", "./out/assets/mp3/各科範圍/"+series.Books[i].Chapters[j].Headings[k].Url+".mp3")
					}

					if doesFileExist("./out/assets/mp3/師父音檔/" + series.Books[i].Chapters[j].Headings[k].Url + ".mp3") {
						log.Info("找到師父音檔!")
					} else {
						log.Errorf("無法找到師父音檔: %v", "./out/assets/mp3/師父音檔/"+series.Books[i].Chapters[j].Headings[k].Url+".mp3")
					}
				}

				if series.Books[i].Chapters[j].Headings[k].Type == "Link" {
					log.Infof("正在檢查連結標題：%v", strings.TrimSpace(series.Books[i].Chapters[j].Headings[k].Name))

					if doesFileExist("./out/" + series.Books[i].Chapters[j].Headings[k].Url + ".html") {
						log.Infof("找到連結: %v", series.Books[i].Chapters[j].Headings[k].Url+".html")
					} else {
						log.Errorf("無法找到連結: %v", series.Books[i].Chapters[j].Headings[k].Url+".html")
					}
				}
			}

		}
	}

}

func saveToFile(filename string, content string) {
	name := "./out/" + filename
	file, err := os.Create(name)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	file.WriteString(content)
	log.Infof("輸出%v", filename)
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

func doesFileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}
