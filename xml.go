package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
)

// Tag general tag of xml
type Tag struct {
	Name     xml.Name
	Attr     []xml.Attr
	Children []interface{}
	From     From
}

// Copy returns copy of Tag
func (t *Tag) Copy() *Tag {
	bs, err := xml.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}
	copy := &Tag{}
	xml.NewDecoder(bytes.NewBuffer(bs)).Decode(copy)
	copy.From = t.From
	return copy
}

func copyChildren(cs []interface{}) []interface{} {
	copy := []interface{}{}
	for _, v := range cs {
		switch v.(type) {
		case *Tag:
			child := v.(*Tag)
			copy = append(copy, child.Copy())
		case xml.CharData:
			copy = append(copy, v.(xml.CharData))
		case xml.Comment:
			copy = append(copy, v.(xml.Comment))
		}
	}
	return copy
}

// MarshalXML to XML from Tag
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

// UnmarshalXML to Tag from XML
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
			//t.Children = append(t.Children, token.(xml.CharData).Copy())
		case xml.Comment:
			//t.Children = append(t.Children, token.(xml.Comment).Copy())
		}
	}
}

// String from Tag
func (t *Tag) String() string {
	buf := new(bytes.Buffer)
	xml.NewEncoder(buf).Encode(t)
	return buf.String()
}
