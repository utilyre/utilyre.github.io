package main

import (
	"html/template"
	"log"
	"os"
)

type Treasure struct {
	Title   string
	Authors []string
	Links   []TreasureLink
	Tags    []string
}

type TreasureLink struct {
	Label string
	URL   string
}

var treasures = []Treasure{
	{
		Title:   "Introduction to Computer Graphics",
		Authors: []string{"Cem Yuksel"},
		Links: []TreasureLink{
			{"YouTube Playlist", "https://youtube.com/playlist?list=PLplnkTzzqsZTfYh4UbhLGpI5kGd5oW_Hh"},
		},
		Tags: []string{"cg", "course"},
	},
	{
		Title: "Ray Tracing in One Weekend",
		Authors: []string{
			"Peter Shierly",
			"Trevor D. Black",
			"Steve Hollasch",
		},
		Links: []TreasureLink{{"Online Book", "https://raytracing.github.io"}},
		Tags:  []string{"cg", "cxx"},
	},
	{
		Title: "A Gentle Introduction to ReSTIR",
		Authors: []string{
			"Chris Wyman",
			"Markus Kettunen",
			"Daqi Lin",
			"Benedikt Bitterli",
			"Cem Yuksel",
			"Wojciech Jarosz",
			"Pawel Kozlowski",
			"Giovanni De Francesco",
		},
		Links: []TreasureLink{{"Online Course", "https://intro-to-restir.cwyman.org"}},
		Tags:  []string{"cg", "course"},
	},
	{
		Title:   "GPU Gems by NVIDIA",
		Authors: []string{"NVIDIA"},
		Links: []TreasureLink{
			{"Volume 1", "https://developer.nvidia.com/gpugems/gpugems"},
			{"Volume 2", "https://developer.nvidia.com/gpugems/gpugems2"},
			{"Volume 3", "https://developer.nvidia.com/gpugems/gpugems3"},
		},
		Tags: []string{"cg"},
	},
	{
		Title:   "Physically Based Rendering",
		Authors: []string{"Matt Pharr", "Wenzel Jakob", "and Greg Humphreys"},
		Links:   []TreasureLink{{"Online Book", "https://pbr-book.org"}},
		Tags:    []string{"cg", "cxx"},
	},
}

func main() {
	t, err := template.ParseFiles("treasure.temp.html")
	if err != nil {
		log.Fatalln(err)
	}

	err = t.ExecuteTemplate(os.Stdout, "treasures", treasures)
	if err != nil {
		log.Fatalln(err)
	}
}
