package main

import (
	"bufio"
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

var sortMap map[string]int

func init() {
	lines := readLines("attrSortOrder.txt")
	sortMap = map[string]int{}
	for i, line := range lines {
		sortMap[line] = i
	}
}

func readLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func convXML(path string, sets *Settings, ccs []string) {
	log.Println("Convert SXML:" + path)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	f, _ := os.Open(path)
	fi, _ := f.Stat()
	fileName := fi.Name()

	root := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)
	conv(root, fileName, sets, ccs)

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

func conv(t *Tag, fileName string, sets *Settings, ccs []string) {
	var tipe, id string

	//log.Println(t.Name.Local)

	for _, a := range t.Attr {
		switch a.Name.Local {
		case "type":
			tipe = a.Value
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
	//log.Println(ss)

	if tipe != "" {
		vmap := make(map[string]Value)
		ids := []string{""}
		if id != "" {
			ids = append(ids, id)
		}
		for _, id := range ids {
			for _, s := range ss {
				if set, ok := (*sets)[tipe+id+s]; ok {
					for k, v := range set.Map {
						vmap[k] = v
					}
					t.Children = append(t.Children, set.Children...)
				}
			}
		}
		for _, a := range t.Attr {
			vmap[a.Name.Local] = Value{a.Value, From{fileName, "!THIS!"}}
		}

		as := []Attr{}
		for k, v := range vmap {
			as = append(as, Attr{Name: k, Value: v})
		}
		sort.Sort(AttrsByName(as))

		if debugFlag {
			need := false
			buf := bytes.NewBufferString("\n")
			if t.From.Name != "" {
				need = true
				fmt.Fprintf(buf, "<%s> from \"%s\" in \"%s\"\n", t.Name.Local, t.From.Selector, t.From.Name)
			}
			for _, a := range as {
				if a.Value.From.Selector != "!THIS!" {
					need = true
					fmt.Fprintf(buf, "%s = \"%s\" from \"%s\" in \"%s\"\n", a.Name, a.Value.Value, a.Value.From.Selector, a.Value.From.Name)
				}
			}
			if need {
				c := []interface{}{xml.Comment(buf.Bytes())}
				t.Children = append(c, t.Children...)
			}
		}

		xas := []xml.Attr{}
		for _, a := range as {
			xas = append(xas, xml.Attr{
				Name:  xml.Name{Space: "", Local: a.Name},
				Value: a.Value.Value,
			})
		}
		t.Attr = xas
	}

	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			conv(tag, fileName, sets, ccs)
		}
	}
}
