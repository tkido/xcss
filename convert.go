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
	lines := readLines("attrsort.txt")
	sortMap = map[string]int{}
	for i, line := range lines {
		sortMap[line] = i
	}
	log.Println(sortMap)
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
	if err2 := scanner.Err(); err2 != nil {
		log.Fatal(err2)
	}
	return lines
}

func convXML(path string, sets *Settings, ccs []string) {
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
	conv(root, fileName, sets, ccs)

	dir := filepath.Dir(path)
	newName := strings.Replace(fileName, "_sxml.xml", ".xml", 1)
	newPath := filepath.Join(dir, newName)

	file, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	xml.NewEncoder(buf).Encode(root)
	str := strings.Replace(buf.String(), "&#x9;", "\t", -1) //Temporary Workaround for a bug of encoder(maybe)
	file.WriteString(str)
}

func conv(t *Tag, fileName string, sets *Settings, ccs []string) {
	var tipe, id string

	log.Println(t.Name.Local)

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
	log.Println(ss)

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
				}
			}
		}
		for _, a := range t.Attr {
			vmap[a.Name.Local] = Value{a.Value, From{fileName, "THIS"}}
		}
		as := []xml.Attr{}
		for k, v := range vmap {
			as = append(as, xml.Attr{
				Name:  xml.Name{Space: "", Local: k},
				Value: v.Value,
			})
		}
		sort.Sort(AttrByName(as))
		t.Attr = as

		if debugFlag {
			buf := bytes.NewBufferString("\n")
			for k, v := range vmap {
				fmt.Fprintf(buf, "%s = \"%s\" from \"%s\" in \"%s\"\n", k, v.Value, v.From.Selector, v.From.Name)
			}

			c := xml.Comment(buf.Bytes())
			t.Children = append(t.Children, c)
		}
	}

	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			conv(tag, fileName, sets, ccs)
		}
	}
}

//AttrByName is []xml.Attr sorted by names in "attrsort.txt"
type AttrByName []xml.Attr

func (p AttrByName) Len() int { return len(p) }
func (p AttrByName) Less(i, j int) bool {
	return sortMap[p[i].Name.Local] < sortMap[p[j].Name.Local]
}
func (p AttrByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
