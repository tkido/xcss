package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	sufSXML = "_sxml.xml"
	sufXCSS = "_xcss.xml"
	sufXML  = ".xml"
)

// ConvSetting is watched SXML file's data
type ConvSetting struct {
	FilePath string
	Settings *Settings
	Updated  time.Time
	/* Due to library limitations, events may be sent twice.
	refer to "Extending Fsnotify" in fsnotify https://fsnotify.org/
	Use this "Updated" value to prevent double conversion. */
}

// WatchSetting is Watcher and watchlist(map) of watched SXML files
type WatchSetting struct {
	Watcher  *fsnotify.Watcher
	WatchMap map[string]*ConvSetting
}

func main() {
	for {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		wset := WatchSetting{watcher, map[string]*ConvSetting{}}
		walk(rootFlag, &Settings{}, wset)

		if !watchFlag {
			break // exit main
		}

		log.Println("Watching...")
	Loop:
		for {
			select {
			case ev := <-watcher.Events:
				reset := false
				if strings.HasSuffix(ev.Name, sufSXML) {
					if ev.Op&fsnotify.Create != 0 {
						reset = true
					} else if ev.Op&fsnotify.Write != 0 {
						if cs, ok := wset.WatchMap[ev.Name]; ok {
							now := time.Now()
							// "Write" event within 0.5 second to the same file is regarded as duplicated.
							if now.Sub(cs.Updated) > time.Second/2 {
								convXML(cs.FilePath, cs.Settings, initClasses)
								cs.Updated = now
							}
						}
					}
				} else if strings.HasSuffix(ev.Name, sufXCSS) {
					if ev.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) != 0 {
						reset = true
					}
				}
				if reset {
					log.Println("Reflesh!!")
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

	var dirs, csss, xmls []os.FileInfo
	for _, fi := range fis {
		if fi.IsDir() {
			dirs = append(dirs, fi)
		} else {
			name := fi.Name()
			if strings.HasSuffix(name, sufXCSS) {
				csss = append(csss, fi)
			} else if strings.HasSuffix(name, sufSXML) {
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
		wset.WatchMap[xmlPath] = &ConvSetting{xmlPath, sets, time.Now()}
	}
	for _, dir := range dirs {
		fullPath := filepath.Join(path, dir.Name())
		walk(fullPath, sets, wset)
	}
	// Register to the watcher last. Because it is not necessary to receive events from the first global conversion.
	err = wset.Watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
}
