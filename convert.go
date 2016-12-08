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

func convXML(path string, sets *Settings) {
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
	conv(root, fileName, sets)

	log.Println(root)
}

func conv(t *Tag, fileName string, sets *Settings) {
	var tipe, id string
	ss := []string{""}

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
				ss = append(ss, comb(cs)...)
			}
		}
		if tipe != "" {
			vmap := make(map[string]Value)
			if id != "" {
				ss = append(ss, id)
			}
			for _, s := range ss {
				if set, ok := (*sets)[tipe+s]; ok {
					for k, v := range set.Map {
						vmap[k] = v
					}
				}
			}
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
			conv(tag, fileName, sets)
		}
	}
}
