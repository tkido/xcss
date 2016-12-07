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
	name := fi.Name()

	root := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)
	conv(root, name)

	log.Println(root)
}

func conv(t *Tag, name string) {
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
				vmap[a.Name.Local] = Value{a.Value, From{name, key}}
			}

			if set, ok := sets[key]; ok {
				for k, v := range set.Map {
					if _, ok := vmap[k]; !ok {
						vmap[k] = v
					}
				}
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
			conv(tag, name)
		}
	}
}
