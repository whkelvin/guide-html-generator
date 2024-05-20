package parser

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Heading struct {
	Name string
	Type string
	Url  string
}
type Chapter struct {
	Filename       string
	Name           string
	Headings       []Heading
	RangeAudioUrl  string
	Next           string
	Prev           string
	TotalBookCount int
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
			s.Books[i].Chapters[j] = *getChapterFromCsv(csvPath, fmt.Sprintf("冊%v表%v", i+1, j+1), bookCount)
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

func getChapterFromCsv(path string, filenameWithoutExtention string, totalBookCount int) *Chapter {
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
	TITLE_COL := 0

	RANGE_AUDIO_PATH_ROW := 1
	RANGE_AUDIO_PATH_COL := 1

	PREV_ROW := 2
	PREV_COL := 1

	NEXT_ROW := 3
	NEXT_COL := 1

	HEADING_COL := 0
	TYPE_COL := 1
	URL_COL := 2

	chapter := Chapter{
		Name:           records[TITLE_ROW][TITLE_COL],
		Filename:       filenameWithoutExtention,
		RangeAudioUrl:  records[RANGE_AUDIO_PATH_ROW][RANGE_AUDIO_PATH_COL],
		Next:           records[NEXT_ROW][NEXT_COL],
		Prev:           records[PREV_ROW][PREV_COL],
		Headings:       []Heading{},
		TotalBookCount: totalBookCount,
	}

	headingCount := len(records) - 5
	if headingCount > 0 {
		for i := 5; i < len(records); i++ {
			chapter.Headings = append(chapter.Headings, Heading{
				Name: records[i][HEADING_COL],
				Type: records[i][TYPE_COL],
				Url:  records[i][URL_COL],
			})
		}
	}
	return &chapter
}
