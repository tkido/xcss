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

func readCSS(path string, sets *Settings) {
	log.Println("Read XCSS:" + path)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Open(path)
	defer f.Close()
	fi, _ := f.Stat()
	fileName := fi.Name()

	root := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)

	// tags just under the root(<styles>) are recognized as styles
	for _, v := range root.Children {
		if tag, isTag := v.(*Tag); isTag {
			parse(tag, fileName, sets)
		}
	}
}

func parse(t *Tag, fileName string, sets *Settings) {
	var key, tipe, id, class string

	as := []xml.Attr{}
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
			default:
				as = append(as, a)
			}
		}
		if tipe != "" {
			key = tipe + id + class
			t.From = From{Name: fileName, Selector: key}

			vmap := make(map[string]Value)
			for _, a := range as {
				vmap[a.Name.Local] = Value{a.Value, From{fileName, key}}
			}

			if set, ok := (*sets)[key]; ok {
				for k, v := range vmap {
					set.Map[k] = v
				}
				// when there is a stronger setting for the same selector, children tags is overwritten
				set.Children = t.Children
			} else {
				(*sets)[key] = &Setting{vmap, t.Children}
			}
		}
	}

}
