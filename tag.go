package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

//From indicates where the value comes from
type From struct {
	Path string
	Tag  *Tag
}

//Value is the value set for an attribute
type Value struct {
	Value string
	From  From
}

//Setting a Set of settings corresponding to one selector
type Setting struct {
	Map      map[string]Value
	Children []interface{}
}

//String
func (set *Setting) String() string {
	return fmt.Sprintf("%v\n", set.Map)
}

//Settings the total settnigs from CSSs in project
type Settings map[string]*Setting

// Tag general tag of xml
type Tag struct {
	Name     xml.Name
	Attr     []xml.Attr
	Children []interface{}
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
			t.Children = append(t.Children, token.(xml.CharData).Copy())
		case xml.Comment:
			t.Children = append(t.Children, token.(xml.Comment).Copy())
		}
	}
}

//String
func (t *Tag) String() string {
	buf := new(bytes.Buffer)
	xml.NewEncoder(buf).Encode(t)
	return buf.String()
}
