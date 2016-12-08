package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var sets Settings

func main() {
	sets = Settings{}
	csss := []string{
		"./testdata/platform/platform_css.xml",
		"./testdata/platform/project/project_css.xml",
	}
	for _, css := range csss {
		readCSS(css)
	}

	appsPath := "./testdata/platform/project/apps"

	convCSS("./testdata/platform/project/apps/foo/foo_main_style.xml")

	err := filepath.Walk(appsPath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				// 特定のディレクトリ以下を無視する場合は
				// return filepath.SkipDir
				return nil
			}
			rel, err := filepath.Rel(appsPath, path)
			fmt.Println(rel)
			return nil
		})

	if err != nil {
		fmt.Println(1, err)
	}

}
