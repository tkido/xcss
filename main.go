package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const (
	SXMLSuffix = "_sxml.xml"
	XCSSSuffix = "_xcss.xml"
	XMLSuffix  = ".xml"
)

// ConvSetting is a pair of watched SXML's filepath and Settings applyed to this file
type ConvSetting struct {
	FilePath string
	Settings *Settings
}

// WatchSetting is a pair of Watcher and list(map) of watched SXML files
type WatchSetting struct {
	Watcher  *fsnotify.Watcher
	WatchMap map[string]ConvSetting
}

func main() {
	for {
		doWalk()
		if !watchFlag {
			break
		}
	}
}

func doWalk() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	wset := WatchSetting{watcher, map[string]ConvSetting{}}
	walk(rootFlag, &Settings{}, wset)

	if watchFlag {
		log.Println("Watching...")
	Loop:
		for {
			select {
			case ev := <-watcher.Events:
				//log.Println(ev)
				reset := false
				if strings.HasSuffix(ev.Name, SXMLSuffix) {
					if ev.Op&fsnotify.Create != 0 {
						reset = true
					} else if ev.Op&fsnotify.Write != 0 {
						if cs, ok := wset.WatchMap[ev.Name]; ok {
							convXML(cs.FilePath, cs.Settings, initClasses)
						}
					}
				} else if strings.HasSuffix(ev.Name, XCSSSuffix) {
					if ev.Op&fsnotify.Chmod == 0 {
						reset = true
					}
				}
				if reset {
					log.Println("Now Reflesh! filesets for settings are modified.")
					break Loop
				}
			case err = <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}
}

func walk(path string, sets *Settings, wset WatchSetting) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	err = wset.Watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	var dirs, csss, xmls []os.FileInfo
	for _, fi := range fis {
		if fi.IsDir() {
			dirs = append(dirs, fi)
		} else {
			name := fi.Name()
			if strings.HasSuffix(name, XCSSSuffix) {
				csss = append(csss, fi)
			} else if strings.HasSuffix(name, SXMLSuffix) {
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
		wset.WatchMap[xmlPath] = ConvSetting{xmlPath, sets}
	}
	for _, dir := range dirs {
		fullPath := filepath.Join(path, dir.Name())
		walk(fullPath, sets, wset)
	}
}
