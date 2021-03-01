package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {

	//unmarshal my resume.json into the structs generated via https://github.com/a-h/generate
	jsonResumeFile, err := os.Open("resume.json")
	if err != nil {
		log.Fatalf("error opening json resume file: %v", err)
	}

	var resume ResumeSchema

	err = json.NewDecoder(jsonResumeFile).Decode(&resume)
	if err != nil {
		log.Fatalf("error decoding json resume file: %v", err)
	}

	err = generateTemplatedFile(resume, "files/author.tmpl", "../data/en/author.yaml")
	if err != nil {
		log.Fatalf("error templating author: %v", err)
	}

	err = generateTemplatedFile(resume, "files/about.tmpl", "../data/en/sections/about.yaml")
	if err != nil {
		log.Fatalf("error templating about: %v", err)
	}

}

func generateTemplatedFile(resume ResumeSchema, srcPath string, destPath string) error {
	// open the templated yaml file
	template, err := template.ParseFiles(srcPath)
	if err != nil {
		return fmt.Errorf("error opening template: %w", err)
	}

	// create new file to be written to
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	// execute template
	err = template.Execute(file, resume)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	file.Close()

	return nil
}
