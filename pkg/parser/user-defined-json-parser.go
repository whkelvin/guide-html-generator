package parser

import (
	"encoding/json"
	. "html_generator/pkg/parser/models"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func readUserGeneratedModelFromJson(filepath string) (*UserGeneratedModel, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var chapter UserGeneratedModel
	err = json.Unmarshal(data, &chapter)
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func plantTree(filename []string, seed string) Tree {
	root := Tree{
		Title:    seed,
		Children: []Tree{},
	}

	childrenStr := getImmediateChildren(filename, seed)
	for i := 0; i < len(childrenStr); i++ {
		root.Children = append(root.Children, plantTree(filename, childrenStr[i]))
	}

	return root
}

type ByNumber []string

func (a ByNumber) Len() int { return len(a) }
func (a ByNumber) Less(i int, j int) bool {

	num1, err := strconv.Atoi(a[i])
	if err != nil {
		panic("not a number")
	}
	num2, err := strconv.Atoi(a[j])
	if err != nil {
		panic("not a number")
	}
	return num1 < num2
}
func (a ByNumber) Swap(i int, j int) { a[i], a[j] = a[j], a[i] }

type BySuffix []string

func (a BySuffix) Len() int { return len(a) }
func (a BySuffix) Less(i int, j int) bool {
	str1 := a[i]
	str2 := a[j]

	index := strings.LastIndex(str1, "-")
	str1 = str1[index+1:]

	index = strings.LastIndex(str2, "-")
	str2 = str2[index+1:]

	num1, err := strconv.Atoi(str1)
	if err != nil {
		panic("not a number")
	}
	num2, err := strconv.Atoi(str2)
	if err != nil {
		panic("not a number")
	}
	return num1 < num2
}

func (a BySuffix) Swap(i int, j int) { a[i], a[j] = a[j], a[i] }

func getImmediateChildren(filenames []string, parent string) []string {
	if parent == "" {
		children := []string{}

		for i := range filenames {
			if !strings.ContainsAny(filenames[i], "-") {
				children = append(children, filenames[i])
			}
		}
		sort.Sort(ByNumber(children))
		return children
	}

	prefix := parent + "-"
	offspring := []string{}

	for i := range filenames {
		if strings.HasPrefix(filenames[i], prefix) {
			offspring = append(offspring, filenames[i])
		}
	}

	numberOfDashesOfParents := strings.Count(parent, "-")
	numberOfDashesOfImmediateChildren := numberOfDashesOfParents + 1
	children := []string{}
	for i := range offspring {
		if strings.Count(offspring[i], "-") == numberOfDashesOfImmediateChildren {
			children = append(children, offspring[i])
		}
	}
	sort.Sort(BySuffix(children))

	return children
}

func readDir(path string) ([]string, error) {
	var files []string
	f, err := os.Open(path)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		name, err := stripExt(file.Name())
		if err != nil {
			name = ""
		}
		files = append(files, name)
	}
	return files, nil
}

func stripExt(filename string) (string, error) {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	return name, nil
}
