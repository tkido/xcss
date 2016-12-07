package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

var settings map[string]*Tag

func main() {
	settings = map[string]*Tag{}
	readCSS("./testdata/platform/platform_css.xml")
	readCSS("./testdata/platform/project/project_css.xml")
}

func readCSS(path string) {
	log.Println("Read CSS:" + path)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	root := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)
	parse(root)
	log.Println(settings)
}

func parse(t *Tag) {
	var key, tipe, id, class string

	log.Println(t.Name.Local)
	if t.Name.Local == "item" {
		for _, a := range t.Attr {
			switch a.Name.Local {
			case "type":
				tipe = a.Value
			case "id":
				id = "#" + a.Value
			case "class":
				cs := strings.Split(a.Value, " ")
				sort.Strings(cs)
				class = "." + strings.Join(cs, ".")
			}
		}
		if tipe != "" {
			key = tipe + id + class
			settings[key] = t
		}
	}

	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			parse(tag)
		}
	}
}
