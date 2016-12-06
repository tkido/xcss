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

	root := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(&root)

	parse(root)
}

func parse(t *Tag) {
	log.Println(t.Name.Local)

	for _, a := range t.Attr {
		log.Println(a.Name.Local)
		log.Println(a.Value)
	}

	for _, v := range t.Children {
		switch v.(type) {
		case *Tag:
			parse(v.(*Tag))
		case xml.CharData:
			//log.Println(string(v.(xml.CharData)))
		case xml.Comment:
			//log.Println(string(v.(xml.Comment)))
		}
	}
}
