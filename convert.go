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

func convCSS(path string) {
	log.Println("Convert CSS:" + path)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	f, _ := os.Open(path)
	fi, _ := f.Stat()
	fileName := fi.Name()

	root := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)
	conv(root, fileName)

	log.Println(root)
}

func apply(vmap map[string]Value, selector string) {
	if set, ok := sets[selector]; ok {
		for k, v := range set.Map {
			vmap[k] = v
		}
	}
}

func conv(t *Tag, fileName string) {
	var tipe, id string
	var css []string //class selectors

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
				css = comb(cs)
			}
		}
		if tipe != "" {
			vmap := make(map[string]Value)
			for _, cs := range css {
				apply(vmap, tipe+cs)
			}
			apply(vmap, tipe+id)
			for _, a := range t.Attr {
				vmap[a.Name.Local] = Value{a.Value, From{fileName, "this"}}
			}
			as := []xml.Attr{}
			for k, v := range vmap {
				as = append(as, xml.Attr{
					Name:  xml.Name{Space: "", Local: k},
					Value: v.Value,
				})
			}
			t.Attr = as
		}
	}

	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			conv(tag, fileName)
		}
	}
}
