package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var attrs = map[string]bool{}

func attrlist(path string) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var xmls []os.FileInfo
	for _, fi := range fis {
		if fi.IsDir() {
		} else {
			xmls = append(xmls, fi)
		}
	}
	for _, xml := range xmls {
		xmlPath := filepath.Join(path, xml.Name())
		readXML(xmlPath)
	}

	for k, _ := range attrs {
		println(k)
	}
}

func readXML(path string) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	root := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)
	parseTag(root)
}

func parseTag(t *Tag) {
	for _, a := range t.Attr {
		attrs[a.Name.Local] = true
	}
	for _, v := range t.Children {
		if tag, isTag := v.(*Tag); isTag {
			parseTag(tag)
		}
	}
}
