package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
)

func main() {
	bs, err := ioutil.ReadFile("./testdata/platform/platform_css.xml")
	if err != nil {
		log.Fatal(err)
	}

	v := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&v)

	log.Println(v.Name)
	log.Println(v.Attr)
	log.Println(v.Children)

}
