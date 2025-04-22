package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const testTemplate = `//go:build {{.Tag}}
// +build {{.Tag}}

package {{.Package}}

import (
	"testing"
)

func Test{{.FuncName | title}}_{{.Tag | title}}(t *testing.T) {
	t.Log("{{.Tag | title}} test for {{.FuncName}}")
}
`

type TemplateData struct {
	Tag      string
	Package  string
	FuncName string
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: testgen <unit|integration> <package> <name> <directory>")
		os.Exit(1)
	}

	tag := os.Args[1]
	packageName := os.Args[2]
	funcName := os.Args[3]
	dir := os.Args[4]

	// Ensure directory exists
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Build full file path
	filename := fmt.Sprintf("%s_test.go", funcName)
	fullPath := filepath.Join(dir, filename)

	// Prepare template
	tmpl := template.Must(template.New("test").Funcs(template.FuncMap{
		"title": func(value string) string {
			return cases.Title(language.English).String(value)
		},
	}).Parse(testTemplate))

	// Create the file
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Execute the template
	err = tmpl.Execute(file, TemplateData{
		Tag:      tag,
		Package:  packageName,
		FuncName: funcName,
	})
	if err != nil {
		fmt.Printf("Error writing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Test file generated at: %s\n", fullPath)
}
