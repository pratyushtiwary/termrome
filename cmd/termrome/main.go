package main

import (
	"encoding/json"
	"fmt"
	"os"

	html_lexer "termrome.io/lexer/html"
)

func Run(filename string) ([]html_lexer.Node, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fmt.Println("Starting to parse", filename)
	parsedContent, err := html_lexer.ParseHtml(string(data))

	if err != nil {
		return nil, err
	}

	return parsedContent, nil
}

func main() {
	if len(os.Args) < 2 {
		panic("Please provide a valid filename")
	}

	var fileName = os.Args[1]

	_, err := os.Stat(fileName)

	if err != nil {
		panic("Please make sure file exists, error: " + err.Error())
	}

	if err != nil {
		fmt.Printf("Lexer error: %s", err.Error())
	}

	parsedContent, err := Run(fileName)

	if err != nil {
		panic("Failed to parse html content, error: " + err.Error())
	}

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
