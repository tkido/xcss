package main

import "log"

var sets Settings

func main() {
	classes := []string{"bar", "foo", "hoge", "piyo"}
	log.Println(comb(classes))

	sets = Settings{}
	readCSS("./testdata/platform/platform_css.xml")
	readCSS("./testdata/platform/project/project_css.xml")
	convCSS("./testdata/platform/project/apps/foo/foo_main.xml")
}
