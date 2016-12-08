package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

var sets Settings

func main() {
	root := "./testdata/platform"
	walk(root, &Settings{})

	/*
		csss := []string{
			"./testdata/platform/platform_css.xml",
			"./testdata/platform/project/project_css.xml",
		}
		for _, css := range csss {
			readCSS(css)
		}

		appsPath := "./testdata/platform/project/apps"

		convCSS("./testdata/platform/project/apps/foo/foo_main_style.xml")

	*/
}

func walk(path string, sets *Settings) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range fis {
		fullPath := filepath.Join(path, fi.Name())

		if fi.IsDir() {
			walk(fullPath, sets)
		} else {
			rel, err := filepath.Rel(path, fullPath)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(rel)
		}
	}
}
