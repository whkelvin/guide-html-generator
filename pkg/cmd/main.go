package main

import (
	"fmt"
	. "html_generator/pkg/generator"
	"html_generator/pkg/parser"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

func main() {
	inputFolderPath := "./assets/input"
	mp3FolderPath := "./assets/mp3"
	src := "./assets/html/wh-audio.js"
	dst := "./out/wh-audio.js"
	_, err := copyFile(src, dst)
	if err != nil {
		panic(err.Error())
	}

	src = "./assets/html/wh-text-audio.js"
	dst = "./out/wh-text-audio.js"
	_, err = copyFile(src, dst)
	if err != nil {
		panic(err.Error())
	}

	seriesToc := parser.GetTOCFromCsv(inputFolderPath + "/index.csv")
	log.Info("找到index.csv...")
	tocItems := []TOCItem{}
	for i := 0; i < len(seriesToc.TOCItems); i++ {
		tocItems = append(tocItems, TOCItem{
			Title: seriesToc.TOCItems[i].Title,
			Type:  seriesToc.TOCItems[i].Type,
			Url:   seriesToc.TOCItems[i].Url,
		})
	}
	seriesTmplInput := SeriesTOCTemplateInput{
		BasePath: mp3FolderPath,
		TOC:      tocItems,
	}
	log.Info("建立根目錄頁...")
	out := GenerateSeriesTOC(seriesTmplInput)
	saveToFile("index.html", out)

	bookFolders := getFolders(inputFolderPath)
	for i := 0; i < len(bookFolders); i++ {
		log.Infof("找到%s資料夾", bookFolders[i].Name())

		tocFilePath := fmt.Sprintf(inputFolderPath+"/%s/index.csv", bookFolders[i].Name())
		log.Infof("正在讀取%s目錄頁: %s ...", bookFolders[i].Name(), tocFilePath)
		bookToc := *parser.GetTOCFromCsv(tocFilePath)

		items := []TOCItem{}
		for n := 0; n < len(bookToc.TOCItems); n++ {
			items = append(items, TOCItem{
				Title: bookToc.TOCItems[n].Title,
				Type:  bookToc.TOCItems[n].Type,
				Url:   bookToc.TOCItems[n].Url,
			})
		}

		bookTemplateInput := BookTOCTemplateInput{
			Items:    items,
			BookName: bookFolders[i].Name(),
			BasePath: mp3FolderPath,
		}
		log.Infof("建立%s目錄頁...", bookFolders[i].Name())
		bookTmplOut := GenerateBookTOC(bookTemplateInput)
		filename := fmt.Sprintf("%s.html", bookFolders[i].Name())
		saveToFile(filename, bookTmplOut)

		bookFolderPath := fmt.Sprintf("%s/%s/", inputFolderPath, bookFolders[i].Name())
		log.Infof("正在讀取%s內容: %s...", bookFolders[i].Name(), bookFolderPath)
		csvs := getFilesExceptForIndexDotCsv(bookFolderPath)

		for j := 0; j < len(csvs); j++ {
			csvPath := bookFolderPath + csvs[j].Name()
			log.Infof("正在讀取%s內容: %s...", bookFolders[i].Name(), csvPath)
			chapters := parser.GetChapterFromCsv(csvPath)
			headings := []Heading{}
			for h := 0; h < len(chapters.Headings); h++ {
				headings = append(headings, Heading{
					Name:        chapters.Headings[h].Name,
					Type:        chapters.Headings[h].Type,
					Url:         chapters.Headings[h].Url,
					LinkBtnText: chapters.Headings[h].BtnName,
				})
			}

			bookNumStr := RemovePrefix(bookFolders[i].Name(), "冊")
			bookNum, err := strconv.Atoi(bookNumStr)
			if err != nil {
				bookNum = 0
			}

			contentTmplInput := ContentTemplateInput{
				BasePath:      mp3FolderPath + "/",
				Headings:      headings,
				Next:          chapters.Next,
				Prev:          chapters.Prev,
				Title:         chapters.Title,
				RangeAudioUrl: chapters.RangeAudioUrl,
				TOCUrl:        bookFolders[i].Name() + ".html",
				BookNum:       bookNum,
			}
			contentTmplOut := GenerateContent(contentTmplInput)
			csvFilenameWithoutSuffix := RemoveSuffix(csvs[j].Name(), ".csv")
			filename := fmt.Sprintf("%s%s.html", bookFolders[i].Name(), csvFilenameWithoutSuffix)
			saveToFile(filename, contentTmplOut)
		}
	}

	log.Infof("網頁輸出完成")
	checkGeneratedSite(inputFolderPath, mp3FolderPath)

}

func checkGeneratedSite(inputFolderPath string, mp3FolderPath string) {
	log.Infof("正在檢查輸出網頁...")
	log.Infof("正在檢查根目錄...")
	if doesFileExist("./out/assets/mp3/操作方式.mp3") == false {
		log.Error("無法找到 ./out/assets/mp3/操作方式.mp3")
	} else {
		log.Info("找到./out/assets/mp3/操作方式.mp3!")
	}

	if doesFileExist("./out/assets/mp3/簡介.mp3") == false {
		log.Error("無法找到 ./out/assets/mp3/簡介.mp3")
	} else {
		log.Info("找到./out/assets/mp3/簡介.mp3!")
	}

	seriesToc := parser.GetTOCFromCsv(inputFolderPath + "/index.csv")
	tocItems := []TOCItem{}
	for i := 0; i < len(seriesToc.TOCItems); i++ {
		tocItems = append(tocItems, TOCItem{
			Title: seriesToc.TOCItems[i].Title,
			Type:  seriesToc.TOCItems[i].Type,
			Url:   seriesToc.TOCItems[i].Url,
		})

		log.Infof("正在檢查根目錄連結: %s", seriesToc.TOCItems[i].Title)
		if seriesToc.TOCItems[i].Type != "Link" {
			log.Errorf("根目錄標題只支援Link")
		}

		if doesFileExist("./out/"+seriesToc.TOCItems[i].Url+".html") == false {
			log.Errorf("無法找到: %s", "./out/"+seriesToc.TOCItems[i].Url+".html")
		} else {
			log.Infof("找到%s", "./out/"+seriesToc.TOCItems[i].Url+".html!")
		}
	}

	bookFolders := getFolders(inputFolderPath)
	for i := 0; i < len(bookFolders); i++ {
		log.Infof("正在檢查%s", bookFolders[i].Name())

		tocFilePath := fmt.Sprintf(inputFolderPath+"/%s/index.csv", bookFolders[i].Name())
		log.Infof("正在檢查%s目錄頁: %s ...", bookFolders[i].Name(), tocFilePath)
		bookToc := *parser.GetTOCFromCsv(tocFilePath)

		for n := 0; n < len(bookToc.TOCItems); n++ {
			if bookToc.TOCItems[n].Type == "Link" {
				if doesFileExist("./out/" + bookToc.TOCItems[n].Url + ".html") {
					log.Infof("找到標題連結: %s | %s | %s", bookFolders[i].Name(), bookToc.TOCItems[n].Title, "./out/"+bookToc.TOCItems[n].Url+".html")
				} else {
					log.Errorf("無法找到標題連結: %s | %s | %s", bookFolders[i].Name(), bookToc.TOCItems[n].Title, "./out/"+bookToc.TOCItems[n].Url+".html")
				}
			}
		}

		bookFolderPath := fmt.Sprintf("%s/%s/", inputFolderPath, bookFolders[i].Name())
		log.Infof("正在檢查%s內容: %s...", bookFolders[i].Name(), bookFolderPath)
		csvs := getFilesExceptForIndexDotCsv(bookFolderPath)

		for j := 0; j < len(csvs); j++ {
			csvPath := bookFolderPath + csvs[j].Name()
			log.Infof("正在檢查%s內容: %s...", bookFolders[i].Name(), csvPath)

			chapter := parser.GetChapterFromCsv(csvPath)

			log.Infof("正在檢查本表範圍音檔: %s", csvPath)
			if doesFileExist("./out/assets/mp3/本表範圍/" + chapter.RangeAudioUrl + ".mp3") {
				log.Info("找到本表範圍音檔!")
			} else {
				log.Errorf("無法找到本表範圍音檔: %s", "./out/assets/mp3/本表範圍/"+chapter.RangeAudioUrl+".mp3")
			}

			log.Infof("正在檢查上一表,下一表,及本冊目錄連結: %s", csvPath)
			if chapter.Prev != "" {
				if doesFileExist("./out/" + chapter.Prev + ".html") {
					log.Info("找到上一表!")
				} else {
					log.Errorf("無法找到上一表: %s", "./out/"+chapter.Prev+".html")
				}
			}

			if chapter.Next != "" {
				if doesFileExist("./out/" + chapter.Next + ".html") {
					log.Info("找到下一表!")
				} else {
					log.Errorf("無法找到下一表: %s", "./out/"+chapter.Next+".html")
				}
			}

			if doesFileExist("./out/" + bookFolders[i].Name() + ".html") {
				log.Infof("找到本冊目錄連結%s!", "./out/"+bookFolders[i].Name()+".html")
			} else {
				log.Errorf("無法找到本冊目錄連結: %s", "./out/"+bookFolders[i].Name()+".html")
			}

			log.Infof("正在檢查標題: %s", csvPath)
			for h := 0; h < len(chapter.Headings); h++ {

				if chapter.Headings[h].Type == "Audio" {
					log.Infof("正在檢查音檔標題: %s", strings.TrimSpace(chapter.Headings[h].Name))

					if doesFileExist("./out/assets/mp3/原文/" + chapter.Headings[h].Url + ".mp3") {
						log.Info("找到原文音檔!")
					} else {
						log.Errorf("無法找到原文音檔: %s", "./out/assets/mp3/原文/"+chapter.Headings[h].Url+".mp3")
					}

					if doesFileExist("./out/assets/mp3/上下層科判/" + chapter.Headings[h].Url + ".mp3") {
						log.Info("找到上下層科判音檔!")
					} else {
						log.Errorf("無法找到上下層科判音檔: %s", "./out/assets/mp3/上下層科判/"+chapter.Headings[h].Url+".mp3")
					}

					if doesFileExist("./out/assets/mp3/各科範圍/" + chapter.Headings[h].Url + ".mp3") {
						log.Info("找到各科範圍音檔!")
					} else {
						log.Errorf("無法找到各科範圍音檔: %s", "./out/assets/mp3/各科範圍/"+chapter.Headings[h].Url+".mp3")
					}

					if doesFileExist("./out/assets/mp3/師父音檔/" + chapter.Headings[h].Url + ".mp3") {
						log.Info("找到師父音檔!")
					} else {
						log.Errorf("無法找到師父音檔: %s", "./out/assets/mp3/師父音檔/"+chapter.Headings[h].Url+".mp3")
					}
				}

				if chapter.Headings[h].Type == "Link" {
					log.Infof("正在檢查連結標題：%s", strings.TrimSpace(chapter.Headings[h].Name))
					if doesFileExist("./out/" + chapter.Headings[h].Url + ".html") {
						log.Infof("找到連結: %s", chapter.Headings[h].Url+".html")
					} else {
						log.Errorf("無法找到連結: %s", chapter.Headings[h].Url+".html")
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
	log.Infof("輸出%s", filename)
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

func countFolders(folderPath string) int {
	d, e := os.ReadDir(folderPath)
	if e != nil {
		return 0
	}
	count := 0
	for i := 0; i < len(d); i++ {
		if d[i].IsDir() {
			count++
		}
	}

	return count
}

func getFolders(path string) []fs.DirEntry {
	d, e := os.ReadDir(path)
	out := []fs.DirEntry{}
	if e != nil {
		return out
	}

	for i := 0; i < len(d); i++ {
		if d[i].IsDir() {
			out = append(out, d[i])
		}
	}

	return out
}

func getFilesExceptForIndexDotCsv(path string) []fs.DirEntry {
	d, e := os.ReadDir(path)
	out := []fs.DirEntry{}
	if e != nil {
		return out
	}

	for i := 0; i < len(d); i++ {
		if d[i].Name() != "index.csv" && strings.HasSuffix(d[i].Name(), ".csv") {
			out = append(out, d[i])
		}
	}

	return out
}

func RemoveSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func RemovePrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):len(s)]
	}
	return s
}
