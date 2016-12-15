package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	//attrlist(rootFlag)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	log.Println("Watching...")
	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				log.Println("event:", ev)
			case err = <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	walk(rootFlag, &Settings{}, watcher)
	log.Println("Waiting signal...")
	s := <-c
	close(c)
	log.Println("signal:", s)
}

func walk(path string, sets *Settings, watcher *fsnotify.Watcher) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add(path)
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
		walk(fullPath, sets, watcher)
	}
}
