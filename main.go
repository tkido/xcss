package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

type Source struct {
	Path string
	Tag  *Tag
}
type Value struct {
	Value  string
	Source Source
}
type Setting struct {
	VMap     map[string]Value
	Children []interface{}
}
type Settings map[string]Setting

var sets Settings

func main() {
	sets = Settings{}
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
	parse(root, path)
	log.Println(sets)
}

func parse(t *Tag, path string) {
	var key, tipe, id, class string
	//src := Source{path, t}

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
			if _, ok := sets[key]; ok {
				sets[key] = Setting{}
			} else {
				sets[key] = Setting{}
			}

		}
	}

	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			parse(tag, path)
		}
	}
}

func attrsToMap(as []xml.Attr) map[string]string {
	m := make(map[string]string)
	for _, a := range as {
		switch a.Name.Local {
		case "type", "id", "class":
			//exclude these attributes
		default:
			m[a.Name.Local] = a.Value
		}
	}
	return m
}
