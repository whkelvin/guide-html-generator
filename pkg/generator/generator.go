package generator

import (
	"fmt"
	"html/template"
	. "html_generator/pkg/constant"
	. "html_generator/pkg/parser/models"
	"os"
)

type TemplateInput struct {
	Chapter     *Chapter
	ParentTitle string
	Navigation  string
	Raw         string
	Translation string
	Recording   string
	AudioType   AudioType
}

type AudioType int64

const (
	Raw         AudioType = 0
	Translation AudioType = 1
	Recording   AudioType = 2
	Navigation  AudioType = 3
)

func GenerateHtml(chapter *Chapter, parentTitle string) error {
	fmt.Println("generating " + chapter.Title + "...")
	fmt.Printf("children count %v \n", len(chapter.Children))

	if len(chapter.Children) == 0 {
		err := generateNodePage(chapter, parentTitle)

		if err != nil {
			return err
		}

		fmt.Println("generated " + chapter.Title + "...")
		return nil
	} else {
		for i := 0; i < len(chapter.Children); i++ {
			fmt.Println("parent title")
			fmt.Println(chapter.Title)
			err := GenerateHtml(&chapter.Children[i], chapter.Title)
			if err != nil {
				return err
			}
		}
		err := generateNavPage(chapter, parentTitle)
		if err != nil {
			return err
		}
		return nil
	}
}

func generateNavPage(chapter *Chapter, parentTitle string) error {
	fmt.Println("reading nav template")
	txt, err := os.ReadFile("assets/html/nav-template.html")
	if err != nil {
		return err
	}

	tmpl, err := template.New("Nav").Parse(string(txt))
	if err != nil {
		return err
	}

	var filename string
	if (chapter.Title) == "" {
		filename = DistPath + "/index.html"
	} else {
		filename = DistPath + "/" + chapter.Title + "-nav" + ".html"
	}

	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		return err
	}

	templateInput := TemplateInput{
		Chapter:     chapter,
		ParentTitle: parentTitle,
		Raw:         chapter.Raw,
		Recording:   chapter.Recording,
		Navigation:  chapter.Navigation,
		Translation: chapter.Translation,
	}
	err = tmpl.Execute(file, templateInput)
	if err != nil {
		return err
	}
	return nil
}

func generateNodePage(chapter *Chapter, parentTitle string) error {
	fmt.Println("reading node template")
	txt, err := os.ReadFile("assets/html/node-template.html")
	if err != nil {
		return err
	}

	tmpl, err := template.New("node").Parse(string(txt))
	if err != nil {
		return err
	}

	filename := DistPath + "/" + chapter.Title + "-node" + ".html"
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		return err
	}

	fmt.Println("executing node template")
	templateInput := TemplateInput{
		Chapter:     chapter,
		ParentTitle: parentTitle,
		Raw:         chapter.Raw,
		Recording:   chapter.Recording,
		Navigation:  chapter.Navigation,
		Translation: chapter.Translation,
	}
	err = tmpl.Execute(file, templateInput)
	if err != nil {
		return err
	}

	if chapter.Raw != "" {
		err = generateAudioPage(chapter, Raw)
		if err != nil {
			return err
		}
	}

	if chapter.Recording != "" {
		err = generateAudioPage(chapter, Recording)
		if err != nil {
			return err
		}
	}

	if chapter.Translation != "" {
		err = generateAudioPage(chapter, Translation)
		if err != nil {
			return err
		}
	}

	if chapter.Navigation != "" {
		err = generateAudioPage(chapter, Navigation)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateAudioPage(chapter *Chapter, audioType AudioType) error {
	fmt.Println("reading audio template")
	txt, err := os.ReadFile("assets/html/audio-template.html")
	if err != nil {
		return err
	}

	tmpl, err := template.New("audio").Parse(string(txt))
	if err != nil {
		return err
	}

	var str string
	switch audioType {
	case Raw:
		str = "-raw"
	case Recording:
		str = "-recording"
	case Translation:
		str = "-translation"
	case Navigation:
		str = "-navigation"
	default:
		str = ""
	}

	filename := DistPath + "/" + chapter.Title + str + "-audio" + ".html"
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	fmt.Println("executing audio template")
	templateInput := TemplateInput{
		Chapter:     chapter,
		ParentTitle: "",
		Raw:         chapter.Raw,
		Recording:   chapter.Recording,
		Navigation:  chapter.Navigation,
		Translation: chapter.Translation,
		AudioType:   audioType,
	}
	err = tmpl.Execute(file, templateInput)
	if err != nil {
		return err
	}

	return nil
}

func Test() error {
	txt, err := os.ReadFile("assets/html/node-template.html")
	if err != nil {
		return err
	}

	tmpl, err := template.New("test").Parse(string(txt))
	if err != nil {
		return err
	}

	var chapter Chapter = Chapter{
		Title:       "hello",
		Raw:         "",
		Recording:   "",
		Translation: "",
		Navigation:  "",
		Children:    []Chapter{},
	}

	err = tmpl.Execute(os.Stdout, chapter)
	if err != nil {
		return err
	}

	//generate parent test
	parenttxt, err := os.ReadFile("assets/html/nav-template.html")
	if err != nil {
		return err
	}

	ptmpl, err := template.New("testParent").Parse(string(parenttxt))
	if err != nil {
		return err
	}

	var pchapter Chapter = Chapter{
		Title:       "hello",
		Raw:         "",
		Recording:   "",
		Translation: "",
		Navigation:  "",
		Children: []Chapter{
			{
				Title:       "hello",
				Raw:         "",
				Recording:   "",
				Translation: "",
				Navigation:  "",
				Children:    []Chapter{},
			},
			{
				Title:       "hello",
				Raw:         "",
				Recording:   "",
				Translation: "",
				Navigation:  "",
				Children:    []Chapter{},
			},
		},
	}

	err = ptmpl.Execute(os.Stdout, pchapter)
	if err != nil {
		return err
	}
	return nil
}
