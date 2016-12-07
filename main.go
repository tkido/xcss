package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
)

var settings map[string]*Tag

func main() {
	settings = map[string]*Tag{}
	readCSS("./testdata/platform/platform_css.xml")
	readCSS("./testdata/platform/project/project_css.xml")
}

func readCSS(path string) {
	log.Println("Read File:" + path)
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
	var key, itemType, selector string

	log.Println(t.Name.Local)
	for _, a := range t.Attr {
		switch a.Name.Local {
		case "type":
			itemType = a.Value
		case "id":
			selector = "#" + a.Value
		case "class":
			selector = "." + a.Value
		}
	}
	if itemType != "" && selector != "" {
		key = itemType + selector
	}
	if key != "" {
		settings[key] = t
	}

	for _, v := range t.Children {
		switch v.(type) {
		case *Tag:
			parse(v.(*Tag))
		}
	}
}
