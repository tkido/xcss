package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	for {
		doWalk()
	}
}

func doWalk() {
	c := make(chan bool, 1)
	defer close(c)

	log.Println("Watching...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Process events
	go func() {
	Loop:
		for {
			select {
			case ev := <-watcher.Events:
				log.Println("event:", ev)
				if ev.Op&fsnotify.Create == fsnotify.Create {
					log.Println("modified file:", ev.Name)
					break Loop
				}
			case err = <-watcher.Errors:
				log.Println("error:", err)
			}
		}
		c <- true
	}()

	walk(rootFlag, &Settings{}, watcher)

	log.Println("Waiting signal...")
	select {
	case b := <-c:
		log.Println(b)
		log.Println("do Walk End")
	}

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
