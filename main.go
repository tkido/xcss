package main

var sets Settings

func main() {
	sets = Settings{}
	readCSS("./testdata/platform/platform_css.xml")
	readCSS("./testdata/platform/project/project_css.xml")
	convCSS("./testdata/platform/project/apps/foo/foo_main.xml")
}
