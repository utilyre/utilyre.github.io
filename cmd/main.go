package main

import (
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
	log.Println("opened content.yml")

	content, err := ingestion.IngestYAML(f)
	if err != nil {
		log.Fatalln("ingest:", err)
	}
	log.Println("ingested yaml contents")

	out, err := os.Create("./public/index.html")
	if err != nil {
		log.Fatalln("create:", err)
	}
	log.Println("created index.html")


	err = rendering.RenderContent(out, content)
	if err != nil {
		log.Fatalln("render:", err)
	}
	log.Println("rendered into index.html")
}
