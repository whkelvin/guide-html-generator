package parser

import (
	. "html_generator/pkg/constant"
	. "html_generator/pkg/parser/models"
)

func Parse() (*Chapter, error) {
	filenames, err := readDir(UserInputJsonPath)
	if err != nil {
		return nil, err
	}
	root := plantTree(filenames, "")

	chapter, err := constructChapterFromTree(root)
	if err != nil {
		return nil, err
	}
	return chapter, nil
}

func constructChapterFromTree(tree Tree) (*Chapter, error) {
	if !tree.HasChild() {
		filename := UserInputJsonPath + "/" + tree.Title + ".json"
		input, err := readUserGeneratedModelFromJson(filename)
		if err != nil {
			return nil, err
		}

		var navigationFilepath = ""
		var translationFilepath = ""
		var recordingFilepath = ""
		var rawFilepath = ""

		if input.Navigation {
			navigationFilepath = NavigationMP3Path + "/" + input.Filename
		}
		if input.Translation {
			translationFilepath = TranslationMP3Path + "/" + input.Filename
		}
		if input.Raw {
			rawFilepath = RawMP3Path + "/" + input.Filename
		}
		if input.Recording {
			recordingFilepath = RecordingMP3Path + "/" + input.Filename
		}

		return &Chapter{
			Title:       input.Title,
			Navigation:  navigationFilepath,
			Recording:   recordingFilepath,
			Raw:         rawFilepath,
			Translation: translationFilepath,
			Children:    []Chapter{},
		}, nil
	}

	var children []Chapter = []Chapter{}
	for i := 0; i < len(tree.Children); i++ {
		chapter, err := constructChapterFromTree(tree.Children[i])
		if err != nil {
			return nil, err
		}
		children = append(children, *chapter)
	}

	filename := UserInputJsonPath + "/" + tree.Title + ".json"
	if tree.IsRoot() {

		return &Chapter{
			Title:       "",
			Raw:         "",
			Translation: "",
			Recording:   "",
			Navigation:  "",
			Children:    children,
		}, nil
	}
	input, err := readUserGeneratedModelFromJson(filename)
	if err != nil {
		return nil, err
	}

	return &Chapter{
		Title:       input.Title,
		Raw:         "",
		Translation: "",
		Recording:   "",
		Navigation:  "",
		Children:    children,
	}, nil
}
