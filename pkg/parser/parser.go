package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"strings"
)

type Heading struct {
	Name    string
	Type    string
	Url     string
	BtnName string
}

type Chapter struct {
	Filename       string
	TOCTitle       string
	Title          string
	Headings       []Heading
	RangeAudioUrl  string
	Next           string
	Prev           string
	TotalBookCount int
	BookNumber     int
}

type Book struct {
	Name     string
	Chapters []Chapter
	TOC      TOC
}

type Series struct {
	Books []Book
}

type TOC struct {
	TOCItems []TOCItem
}

type TOCItem struct {
	Title string
	Type  string
	Url   string
}

func ParseFolder(path string) *Series {
	log.Info("Parsing Folder ...")
	bookCount := countItems(path)
	s := Series{
		Books: make([]Book, bookCount),
	}
	log.Infof("%v books found!", bookCount)

	for i := 0; i < bookCount; i++ {
		s.Books[i].Name = fmt.Sprintf("冊%v", i+1)
		files := listFiles(path + "/" + s.Books[i].Name)

		s.Books[i].Chapters = []Chapter{}

		tocFilePath := fmt.Sprintf("%v/冊%v/index.csv", path, i+1)
		s.Books[i].TOC = *getTOCFromCsv(tocFilePath)

		for j := 0; j < len(files); j++ {
			if files[j] != "index.csv" {
				csvPath := fmt.Sprintf("%v/冊%v/%v", path, i+1, files[j])
				filenameWithoutExtention := RemoveSuffix(files[j], ".csv")
				log.Infof("Getting Chapter Infomation: %v : %v", s.Books[i].Name, filenameWithoutExtention)
				chapter := *getChapterFromCsv(csvPath, fmt.Sprintf("冊%v%v", i+1, filenameWithoutExtention), bookCount, i+1)
				s.Books[i].Chapters = append(s.Books[i].Chapters, chapter)
			}
		}
		log.Infof("%v chapters processed for this book.", len(s.Books[i].Chapters))
	}
	return &s
}

func countItems(folderPath string) int {
	d, e := os.ReadDir(folderPath)
	if e != nil {
		return 0
	}
	return len(d)
}

func RemoveSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func listFiles(folderPath string) []string {
	var filenames []string

	// Open the directory
	dir, err := os.Open(folderPath)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	// Read the directory's contents
	files, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}

	// Loop through the files and collect their names
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}

	return filenames
}

func getTOCFromCsv(path string) *TOC {
	file, err := os.Open(path)
	// Checks for the error
	if err != nil {
		panic("Error while reading the file: " + path + " " + err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Error reading records")
	}

	TOC_TITLE_COl := 0
	TOC_TYPE_COl := 1
	TOC_URL_COl := 2

	log.Info("Parsing TOC...")
	itemCount := len(records) - 1 // not counting the header row
	tocItems := make([]TOCItem, itemCount)

	toc := TOC{
		TOCItems: tocItems,
	}

	for i := 0; i < itemCount; i++ {
		log.Infof("Processing TOC title %v", records[i+1][TOC_TITLE_COl])
		toc.TOCItems[i].Title = records[i+1][TOC_TITLE_COl]
		toc.TOCItems[i].Type = records[i+1][TOC_TYPE_COl]
		toc.TOCItems[i].Url = records[i+1][TOC_URL_COl]
	}

	return &toc
}

func getChapterFromCsv(path string, filenameWithoutExtention string, totalBookCount int, bookNumber int) *Chapter {
	file, err := os.Open(path)
	// Checks for the error
	if err != nil {
		panic("Error while reading the file: " + path + " " + err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Error reading records")
	}

	TOC_TITLE_ROW := 0
	TOC_TITLE_COL := 1

	TITLE_ROW := 0
	TITLE_COL := 3

	RANGE_AUDIO_PATH_ROW := 1
	RANGE_AUDIO_PATH_COL := 1

	PREV_ROW := 2
	PREV_COL := 1

	NEXT_ROW := 3
	NEXT_COL := 1

	HEADING_COL := 0
	TYPE_COL := 1
	URL_COL := 2

	LINK_BTN_NAME_COL := 3

	chapter := Chapter{
		TOCTitle:       records[TOC_TITLE_ROW][TOC_TITLE_COL],
		Title:          records[TITLE_ROW][TITLE_COL],
		Filename:       filenameWithoutExtention,
		RangeAudioUrl:  records[RANGE_AUDIO_PATH_ROW][RANGE_AUDIO_PATH_COL],
		Next:           records[NEXT_ROW][NEXT_COL],
		Prev:           records[PREV_ROW][PREV_COL],
		Headings:       []Heading{},
		TotalBookCount: totalBookCount,
		BookNumber:     bookNumber,
	}

	headingCount := len(records) - 5
	if headingCount > 0 {
		for i := 5; i < len(records); i++ {
			chapter.Headings = append(chapter.Headings, Heading{
				Name:    records[i][HEADING_COL],
				Type:    records[i][TYPE_COL],
				Url:     records[i][URL_COL],
				BtnName: records[i][LINK_BTN_NAME_COL],
			})
		}
	}

	return &chapter
}
