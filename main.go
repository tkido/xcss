package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	walk(rootFlag, &Settings{})
	//attrlist(rootFlag)
}

func walk(path string, sets *Settings) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var dirs, csss, xmls []os.FileInfo
	for _, fi := range fis {
		if fi.IsDir() {
			dirs = append(dirs, fi)
		} else {
			name := fi.Name()
			if strings.HasSuffix(name, "_xcss.xml") {
				csss = append(csss, fi)
			} else if strings.HasSuffix(name, "_sxml.xml") {
				xmls = append(xmls, fi)
			}
		}
	}
	if 0 < len(csss) {
		sets = sets.Copy()
		for _, css := range csss {
			cssPath := filepath.Join(path, css.Name())
			readCSS(cssPath, sets)
		}
	}
	for _, xml := range xmls {
		xmlPath := filepath.Join(path, xml.Name())
		convXML(xmlPath, sets, initClasses)
	}
	for _, dir := range dirs {
		fullPath := filepath.Join(path, dir.Name())
		walk(fullPath, sets)
	}
}
