package main

import (
	"fmt"
	"log"
	"os"
	"website/internal/ingestion"
	"website/internal/rendering"
)

func main() {
	f, err := os.Open("./content.yml")
	if err != nil {
		log.Fatalln("open:", err)
	}

	content, err := ingestion.IngestYAML(f)
	if err != nil {
		log.Fatalln("ingest:", err)
	}

	fmt.Println(len(content.Journal))

	out, err := os.Create("./public/index.html")
	if err != nil {
		log.Fatalln("create:", err)
	}

	fmt.Println(out)

	err = rendering.RenderContent(out, content)
	if err != nil {
		log.Fatalln("render:", err)
	}
}
