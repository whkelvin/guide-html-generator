package main

import (
	"fmt"
	. "html_generator/pkg/generator"
)

func main() {
	templateInput := BooksTOCTemplateInput{
		BookNames: []string{
			"冊1",
			"冊2",
		},
	}
	out := GenerateBooksTOC(templateInput)
	fmt.Println(out)
}
