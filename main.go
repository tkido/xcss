package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	debugFlag bool
	rootFlag  string
	classFlag string

	initClasses []string

	reCSS = regexp.MustCompile(`_css.xml$`)
	reXML = regexp.MustCompile(`_style.xml$`)
	reTab = regexp.MustCompile(`&#x9;`)
)

func init() {
	flag.BoolVar(&debugFlag, "debug", false, "debug flag")
	flag.BoolVar(&debugFlag, "d", false, "debug flag")
	flag.StringVar(&rootFlag, "root", "", "root flag")
	flag.StringVar(&rootFlag, "r", "", "root flag")
	flag.StringVar(&classFlag, "class", "", "class flag")
	flag.StringVar(&classFlag, "c", "", "class flag")
	flag.Parse()

	initClasses = []string{}
	if len(classFlag) > 0 {
		initClasses = strings.Split(classFlag, " ")
	}

}

func main() {
	walk("./testdata/platform", &Settings{})
	log.Println(debugFlag)
	log.Println(rootFlag)
	log.Println(classFlag)
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
			if reCSS.MatchString(name) {
				csss = append(csss, fi)
			} else if reXML.MatchString(name) {
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
		convXML(xmlPath, sets)
	}
	for _, dir := range dirs {
		fullPath := filepath.Join(path, dir.Name())
		walk(fullPath, sets)
	}
}
