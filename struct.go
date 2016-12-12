package main

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"log"
)

// From indicates where the value comes from
type From struct {
	Name     string
	Selector string
}

// Value is the value set for an attribute
type Value struct {
	Value string
	From  From
}

// Setting a Set of settings corresponding to one selector
type Setting struct {
	Map      map[string]Value
	Children []interface{}
}

// String from Setting
func (set *Setting) String() string {
	return fmt.Sprintf("%v\n", set.Map)
}

//Settings is the settings from XCSSs
type Settings map[string]*Setting

// Copy returns deep copy of Settings
func (sets *Settings) Copy() *Settings {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	err := enc.Encode(sets)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	var copy Settings
	err = dec.Decode(&copy)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return &copy
}

//Attr is
type Attr struct {
	Name  string
	Value Value
}

//AttrsByName is []Attr sorted by names in "attrsort.txt"
type AttrsByName []Attr

func (p AttrsByName) Len() int { return len(p) }
func (p AttrsByName) Less(i, j int) bool {
	return sortMap[p[i].Name] < sortMap[p[j].Name]
}
func (p AttrsByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

//AttrByName is []xml.Attr sorted by names in "attrsort.txt"
type AttrByName []xml.Attr

func (p AttrByName) Len() int { return len(p) }
func (p AttrByName) Less(i, j int) bool {
	return sortMap[p[i].Name.Local] < sortMap[p[j].Name.Local]
}
func (p AttrByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
