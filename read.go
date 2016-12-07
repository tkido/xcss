package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func readCSS(path string) {
	log.Println("Read CSS:" + path)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	f, _ := os.Open(path)
	fi, _ := f.Stat()
	name := fi.Name()

	root := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)
	parse(root, name)
	log.Println(sets)
}

func parse(t *Tag, path string) {
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
			vmap := make(map[string]Value)
			for _, a := range t.Attr {
				switch a.Name.Local {
				case "type", "id", "class":
					//exclude these attributes
				default:
					vmap[a.Name.Local] = Value{a.Value, From{path, key}}
				}
			}

			if set, ok := sets[key]; ok {
				for k, v := range vmap {
					set.Map[k] = v
				}
				set.Children = t.Children
			} else {
				sets[key] = &Setting{vmap, t.Children}
			}
		}
	}

	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			parse(tag, path)
		}
	}
}
