package main

import (
	"bytes"
	"encoding/gob"
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
	gob.Register(*root)

	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)

	from := From{fileName, ""}
	parse(root, from, sets, 0)

	// log.Println(sets)
}

func parse(t *Tag, from From, sets *Settings, depth int) {
	var key, tipe, id, class string

	// tags just under the root(<styles>) are recognized as styles
	if depth == 1 {
		as := []xml.Attr{}
		for _, a := range t.Attr {
			switch a.Name.Local {
			case "type":
				tipe = ":" + a.Value
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

		key = t.Name.Local + tipe + id + class
		from.Selector = key

		amap := make(map[string]Value)
		for _, a := range as {
			amap[a.Name.Local] = Value{a.Value, from}
		}

		if set, ok := (*sets)[key]; ok {
			for k, v := range amap {
				set.Map[k] = v
			}
			// when there is a stronger setting for the same selector, children tags is overwritten
			set.Children = t.Children
		} else {
			(*sets)[key] = &Setting{amap, t.Children}
		}
	}
	// in other depth, only set "From" of child elements for debug
	t.From = from
	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			parse(tag, from, sets, depth+1)
		}
	}

}
