package main

import "fmt"

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
