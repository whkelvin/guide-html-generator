package parser

import (
	"encoding/csv"
	"github.com/charmbracelet/log"
	"os"
)

type TOC struct {
	TOCItems []TOCItem
}

type TOCItem struct {
	Title string
	Type  string
	Url   string
}

func GetTOCFromCsv(path string) *TOC {
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

	itemCount := len(records) - 1 // not counting the header row
	tocItems := make([]TOCItem, itemCount)

	toc := TOC{
		TOCItems: tocItems,
	}

	for i := 0; i < itemCount; i++ {
		log.Infof("找到目錄頁標題: %v", records[i+1][TOC_TITLE_COl])
		toc.TOCItems[i].Title = records[i+1][TOC_TITLE_COl]
		toc.TOCItems[i].Type = records[i+1][TOC_TYPE_COl]
		toc.TOCItems[i].Url = records[i+1][TOC_URL_COl]
	}

	return &toc
}

type Heading struct {
	Name    string
	Type    string
	Url     string
	BtnName string
}

type Chapter struct {
	//Filename       string
	Title         string
	Headings      []Heading
	RangeAudioUrl string
	Next          string
	Prev          string
	//TotalBookCount int
	//BookNumber     int
}

func GetChapterFromCsv(path string) *Chapter {
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

	TITLE_ROW := 0
	TITLE_COL := 1

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
		Title:         records[TITLE_ROW][TITLE_COL],
		RangeAudioUrl: records[RANGE_AUDIO_PATH_ROW][RANGE_AUDIO_PATH_COL],
		Next:          records[NEXT_ROW][NEXT_COL],
		Prev:          records[PREV_ROW][PREV_COL],
		Headings:      []Heading{},
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
