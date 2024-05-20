package parser

import (
	"encoding/csv"
	"fmt"
	"os"
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
}
type Series struct {
	Books []Book
}

func ParseFolder(path string) *Series {
	bookCount := countItems(path)
	s := Series{
		Books: make([]Book, bookCount),
	}

	for i := 0; i < bookCount; i++ {
		s.Books[i].Name = fmt.Sprintf("冊%v", i+1)
		chapterCount := countItems(path + "/" + s.Books[i].Name)
		s.Books[i].Chapters = make([]Chapter, chapterCount)
		for j := 0; j < chapterCount; j++ {
			csvPath := fmt.Sprintf("%v/冊%v/表%v.csv", path, i+1, j+1)
			s.Books[i].Chapters[j] = *getChapterFromCsv(csvPath, fmt.Sprintf("冊%v表%v", i+1, j+1), bookCount, i+1)
		}
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
