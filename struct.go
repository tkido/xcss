package main

import (
	"bytes"
	"encoding/gob"
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
	return fmt.Sprintf("%+v\n%+v\n", set.Map, set.Children)
}

// Copy returns copy of Setting
func (set *Setting) Copy() *Setting {
	copy := Setting{}
	// child elements of one selector may be replaced by same selector's one.
	// but the element itself is never changed, so there is no problem with shallow copy
	copy.Children = set.Children

	// make a deep copy of Setting.Map
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	err := enc.Encode(set.Map)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	err = dec.Decode(&(copy.Map))
	if err != nil {
		log.Fatal("decode error:", err)
	}

	return &copy
}

//Settings is the settings from XCSSs
type Settings map[string]*Setting

// Copy returns copy of Settings
func (sets *Settings) Copy() *Settings {
	copy := Settings{}
	for k, v := range *sets {
		copy[k] = v.Copy()
	}
	return &copy
}

//Attr is pair of Name, Value. but Value includes Form
type Attr struct {
	Name  string
	Value Value
}
