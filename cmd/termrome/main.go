package main

import (
	"encoding/json"
	"fmt"
	"os"

	"termrome.io/lexer"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please provide a valid filename")
	}

	var fileName string = os.Args[1]

	_, err := os.Stat(fileName)

	if err != nil {
		panic("Please make sure file exists, error: " + err.Error())
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting to parse", fileName)
	parsedContent := lexer.ParseHtml(string(data))

	fmt.Println("File parsed successfully")

	jsonData, err := json.MarshalIndent(parsedContent, "", "	")

	if err != nil {
		panic("Failed to convert AST to JSON, error: " + err.Error())
	}

	fmt.Println("Saving AST in a JSON file...")

	err = os.WriteFile("output.json", jsonData, 0644)

	if err != nil {
		panic("Failed to save file, error: " + err.Error())
	}

	fmt.Println("Output saved successfully to output.json")
}
