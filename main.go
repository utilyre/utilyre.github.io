package main

import (
	"bytes"
	"fmt"
	"html/template"
	textemplate "text/template"
	"io"
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

func main() {
	content, err := readContent()
	if err != nil {
		log.Fatalln(err)
	}

	Skeleton, err := textemplate.ParseFiles("./templates/Skeleton.html")
	if err != nil {
		log.Fatalln(err)
	}

	ResourceList, err := template.ParseFiles("./templates/ResourceList.html")
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	err = ResourceList.Execute(&buf, content.Resources)
	if err != nil {
		log.Fatalln(err)
	}

	err = Skeleton.Execute(os.Stdout, buf.String())
	if err != nil {
		log.Fatalln(err)
	}
}

type Content struct {
	Resources []Resource
}

type Resource struct {
	Title        string
	Associations []string
	Authors      []string
	Links        []string
	Kind         string
	Tags         []string
}

func readContent() (Content, error) {
	f, err := os.Open("./content.yml")
	if err != nil {
		return Content{}, fmt.Errorf("os.Open: %v", err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return Content{}, fmt.Errorf("io.ReadAll: %v", err)
	}

	var content Content
	err = yaml.Unmarshal(data, &content)
	if err != nil {
		return Content{}, fmt.Errorf("yaml.Unmarshal: %v", err)
	}

	return content, nil
}
