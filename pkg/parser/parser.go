package parser

import (
	"os"
	"strings"
)

type Series struct {
	Name  string
	TOC   TOC
	Books []Book
}

type Book struct {
	Name     string
	Chapters []Chapter
	TOC      TOC
}

//func ParseFolder(path string) *Series {
//	log.Info("開始掃描input資料夾 ...")
//	log.Info("開始掃描輯目錄 ...")
//	seriesToc := *GetTOCFromCsv(path + "/index.csv")
//
//	bookCount := countItems(path) - 1 // exclude index.csv
//	s := Series{
//		Books: make([]Book, bookCount),
//		TOC:   seriesToc,
//	}
//	log.Infof("找到%v冊內容", bookCount)
//
//	for i := 0; i < bookCount+1; i++ {
//		s.Books[i].Name = fmt.Sprintf("冊%v", i+1)
//		files := listFiles(path + "/" + s.Books[i].Name)
//
//		s.Books[i].Chapters = []Chapter{}
//
//		tocFilePath := fmt.Sprintf("%v/冊%v/index.csv", path, i+1)
//		log.Infof("正在讀取冊%v目錄頁...", i+1)
//		s.Books[i].TOC = *GetTOCFromCsv(tocFilePath)
//
//		for j := 0; j < len(files); j++ {
//			if files[j] != "index.csv" {
//				csvPath := fmt.Sprintf("%v/冊%v/%v", path, i+1, files[j])
//				filenameWithoutExtention := RemoveSuffix(files[j], ".csv")
//				log.Infof("正在讀取 %v : %v 內容", s.Books[i].Name, filenameWithoutExtention)
//				chapter := *getChapterFromCsv(csvPath, fmt.Sprintf("冊%v%v", i+1, filenameWithoutExtention), bookCount, i+1)
//				s.Books[i].Chapters = append(s.Books[i].Chapters, chapter)
//			}
//		}
//		log.Infof("讀取共%v表", len(s.Books[i].Chapters))
//	}
//	return &s
//}

//func countItems(folderPath string) int {
//	d, e := os.ReadDir(folderPath)
//	if e != nil {
//		return 0
//	}
//	return len(d)
//}

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
