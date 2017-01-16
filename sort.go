package main

import "strings"

var sortOrder map[string]int

func init() {
	lines := strings.Split(strings.TrimSpace(attributesList), "\n")
	sortOrder = map[string]int{}
	for i, line := range lines {
		sortOrder[line] = i
	}
}

// Enumerate all attributes that may appear in layout XMLs.
// Attributes are sorted in the order they appear in this list.
// If there are attributes that do not exist in this list, it logs warning.
var attributesList = `
widget
accessor
value0
type
id
class
x
y
w
h
src-x
src-y
src-w
src-h
leftdim
leftsize
rightdim
rightsize
topdim
topsize
bottomdim
bottomsize
left
right
top
bottom
reducedx
reducedy
min-w
min-h
max-w
max-h
max-padding
itemwidth
itemheight
fg
bg
fontsize
fontcolor
fontstyle
ref
side
size
itemfill
number
layout
direction
scrollrate
alignment
vert-alignment
autoarrange
hinttext
hintsize
hintcolor
lockcolor
maxchar
bgimage
text
url
imageurl
user
touch
disableviewscroll
fill
visiblerows
visiblecolumns
userstates
emptydraw
steps
imagecolumn
isloop
delay
value
`

//AttrsByList is []Attr sorted by names in "attributesList" value
type AttrsByList []Attr

func (p AttrsByList) Len() int { return len(p) }
func (p AttrsByList) Less(i, j int) bool {
	return sortOrder[p[i].Name] < sortOrder[p[j].Name]
}
func (p AttrsByList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
