package main

import (
	"encoding/json"
	"fmt"
	"html_generator/pkg/generator"
	"html_generator/pkg/parser"
	"os"
)

func main() {
	chapter, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
	}

	byte, err := json.Marshal(chapter)
	if err := os.WriteFile("./debug/generated.json", byte, 0666); err != nil {
		fmt.Println("cannot convert to json")
	}

	err = generator.GenerateHtml(chapter, "index")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Done")
}
