package main

import (
	"log"
	"sort"
)

var sets Settings

func main() {
	classes := []string{"bar", "foo", "hoge"}
	log.Println(comb(classes))

	sets = Settings{}
	readCSS("./testdata/platform/platform_css.xml")
	readCSS("./testdata/platform/project/project_css.xml")
	convCSS("./testdata/platform/project/apps/foo/foo_main.xml")
}

func comb(classes []string) [][]string {
	n := len(classes)
	count := 1 << uint(n)

	bits := make(BitsArray, count)
	for i := 0; i < count; i++ {
		bits[i] = i
		log.Printf("%b\n", i)
	}
	sort.Sort(bits)

	aa := [][]string{}
	for i := 0; i < count; i++ {
		log.Printf("%b\n", bits[i])
		a := []string{}
		for j := 0; j < n; j++ {
			if (1<<uint(j))&bits[i] != 0 {
				a = append(a, classes[j])
			}
		}
		aa = append(aa, a)
	}
	return aa
}

//BitsArray sorted by bit count
type BitsArray []int

func (p BitsArray) Len() int { return len(p) }
func (p BitsArray) Less(i, j int) bool {
	if numOfBits(p[i]) == numOfBits(p[j]) {
		return p[i] > p[j]
	}
	return numOfBits(p[i]) > numOfBits(p[j])
}
func (p BitsArray) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func numOfBits(bits int) int {
	bits = (bits & 0x55555555) + (bits >> 1 & 0x55555555)
	bits = (bits & 0x33333333) + (bits >> 2 & 0x33333333)
	bits = (bits & 0x0f0f0f0f) + (bits >> 4 & 0x0f0f0f0f)
	bits = (bits & 0x00ff00ff) + (bits >> 8 & 0x00ff00ff)
	return (bits & 0x0000ffff) + (bits >> 16 & 0x0000ffff)
}
