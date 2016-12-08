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
	return fmt.Sprintf("%v\n", set.Map)
}

//Settings the total settnigs from CSSs in project
type Settings map[string]*Setting

// Copy Settings
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
