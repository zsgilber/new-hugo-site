package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {

	// I want to use my programming resume view for this site:
	programmingResumePath := "views/programming_resume.json"

	//unmarshal my resume.json into the structs generated via https://github.com/a-h/generate
	jsonResumeFile, err := os.Open(programmingResumePath)
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

	groupedWorkItems := groupWorkItemsBy(resume, func(workItems WorkItems) string {
		return workItems.Name
	})

	for _, groupedItem := range groupedWorkItems {
		fmt.Printf("key: %v\n", groupedItem.Key)
	}

	err = generateTemplatedFile(groupedWorkItems, "files/experiences.tmpl", "../data/en/sections/experiences.yaml")
	if err != nil {
		log.Fatalf("error templating about: %v", err)
	}

}

func generateTemplatedFile(data interface{}, srcPath string, destPath string) error {
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
	err = template.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	file.Close()

	return nil
}

// This easier to do here than in the template, the Toha yaml expects positions grouped by company, but the
// JSON resume just has a list, go let's do a groupBy. This need to be sorted also, so make it a slice.
func groupWorkItemsBy(resume ResumeSchema, keySelector func(workItems WorkItems) string) []GroupedWorkItem {
	groupBy := make([]GroupedWorkItem, 0)
	for _, workItem := range resume.Work {
		key := keySelector(*workItem)

		if keyIndex, ok := keyIsInGroupBy(key, groupBy); ok {
			groupBy[keyIndex].Items = append(groupBy[keyIndex].Items, *workItem)
		} else {
			itemToAppend := GroupedWorkItem{
				Key:   key,
				Items: []WorkItems{*workItem},
			}
			groupBy = append(groupBy, itemToAppend)
		}
	}
	return groupBy
}

func keyIsInGroupBy(key string, groupBy []GroupedWorkItem) (index int, ok bool) {
	for i, item := range groupBy {
		if item.Key == key {
			return i, true
		}
	}
	return -1, false
}

type GroupedWorkItem struct {
	Key   string
	Items []WorkItems
}
