package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"

	"github.com/tkido/gotools/log"
)

type Tag struct {
	Name     xml.Name
	Attr     []xml.Attr
	Children []interface{}
}

func (t *Tag) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = t.Name
	start.Attr = t.Attr
	e.EncodeToken(start)
	for _, v := range t.Children {
		switch v.(type) {
		case *Tag:
			child := v.(*Tag)
			if err := e.Encode(child); err != nil {
				return err
			}
		case xml.CharData:
			e.EncodeToken(v.(xml.CharData))
		case xml.Comment:
			e.EncodeToken(v.(xml.Comment))
		}
	}
	e.EncodeToken(start.End())
	return nil
}

func (t *Tag) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t.Name = start.Name
	t.Attr = start.Attr
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			var data *Tag
			if err := d.DecodeElement(&data, &tok); err != nil {
				return err
			}
			t.Children = append(t.Children, data)
		case xml.CharData:
			t.Children = append(t.Children, token.(xml.CharData).Copy())
		case xml.Comment:
			t.Children = append(t.Children, token.(xml.Comment).Copy())
		}
	}
}

const s = `<styles>
	<item type="image" id="background" w="800" h="416" imageurl="@DNA/common/bg.png"/>
</styles>
`

func main() {
	defer log.Close()

	v := &Tag{}
	xml.NewDecoder(bytes.NewBuffer([]byte(s))).Decode(&v)

	log.D(v.Name)
	log.D(v.Attr)
	log.D(v.Children)

	log.D(xml.NewEncoder(os.Stdout).Encode(v))
}
