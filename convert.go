package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func convXML(path string, sets *Settings, ccs []string) {
	log.Println("Convert SXML:" + path)
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
	conv(root, sets, ccs)

	dir := filepath.Dir(path)
	newName := strings.Replace(fileName, sufSXML, sufXML, 1)
	newPath := filepath.Join(dir, newName)
	log.Println("      to XML:" + newPath)

	file, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	output, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	file.Write(output)
}

func conv(t *Tag, sets *Settings, ccs []string) {
	var tipe, id string

	for _, a := range t.Attr {
		switch a.Name.Local {
		case "type":
			tipe = ":" + a.Value
		case "id":
			id = "#" + a.Value
		case "class":
			dup := make([]string, len(ccs))
			copy(dup, ccs)
			ccs = dup

			cs := strings.Split(a.Value, " ")
			ccs = append(ccs, cs...)
		}
	}
	ss := comb(ccs)

	vmap := make(map[string]Value)
	ids := []string{""}
	if id != "" {
		ids = append(ids, id)
	}
	for _, id := range ids {
		for _, s := range ss {
			if set, ok := (*sets)[t.Name.Local+tipe+id+s]; ok {
				for k, v := range set.Map {
					vmap[k] = v
				}
				// when tag's setting applies to multiple selectors,
				// attributes are overwritten by stronger selectors,
				// but children are appended
				t.Children = append(set.Children, t.Children...)
			}
		}
	}

	from := From{"", ""}
	if t.From.Selector != "" {
		from = t.From
	}
	for _, a := range t.Attr {
		vmap[a.Name.Local] = Value{a.Value, from}
	}
	// convert to []Attr from map and sort
	as := []Attr{}
	for k, v := range vmap {
		as = append(as, Attr{Name: k, Value: v})
	}
	sort.Sort(AttrsByList(as))
	// check typo
	for _, a := range as {
		_, ok := sortOrder[a.Name]
		if !ok {
			log.Printf("WARNING: \"%s\" is an unfamiliar attribute. Is not it typo?\n", a.Name)
		}
	}
	// output debug comments
	if flags.Debug {
		need := false
		buf := bytes.NewBufferString("\n")
		if t.From.Name != "" {
			need = true
			fmt.Fprintf(buf, "<%s> from \"%s\" in \"%s\"\n", t.Name.Local, t.From.Selector, t.From.Name)
		}
		for _, a := range as {
			if a.Value.From.Selector != "" {
				need = true
				fmt.Fprintf(buf, "%s = \"%s\" from \"%s\" in \"%s\"\n", a.Name, a.Value.Value, a.Value.From.Selector, a.Value.From.Name)
			}
		}
		if need {
			c := []interface{}{xml.Comment(buf.Bytes())}
			t.Children = append(c, t.Children...)
		}
	}
	// convert to xml.Attr from Attr
	xas := []xml.Attr{}
	for _, a := range as {
		xas = append(xas, xml.Attr{
			Name:  xml.Name{Space: "", Local: a.Name},
			Value: a.Value.Value,
		})
	}
	t.Attr = xas
	// apply recursively to child tags
	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			conv(tag, sets, ccs)
		}
	}
}
