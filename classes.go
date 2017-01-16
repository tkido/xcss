package main

import (
	"log"
	"sort"
	"strings"
)

func comb(classes []string) []string {
	// delete duplicate classes
	dists := make([]string, 0, len(classes))
	checked := map[string]bool{}
	for _, c := range classes {
		if checked[c] {
			log.Printf("WARNING: \"%s\" class is duplicated! It may be unintended.\n", c)
		} else {
			checked[c] = true
			dists = append(dists, c)
		}
	}
	classes = dists

	sort.Sort(sort.Reverse(sort.StringSlice(classes)))
	return getComb(classes)
}

// cache of bit pattern
var bitsMemo = map[int][]int{}

func getBits(n int) []int {
	if cache, ok := bitsMemo[n]; ok {
		return cache
	}
	count := 1 << uint(n)
	bits := make([]int, count)
	for i := 0; i < count; i++ {
		bits[i] = i
	}
	sort.Sort(IntByBits(bits))
	bitsMemo[n] = bits
	return bits
}

// cache of combination of classes
var combMemo = map[string][]string{}

func getComb(classes []string) []string {
	key := strings.Join(classes, ".")
	if cache, ok := combMemo[key]; ok {
		return cache
	}

	n := len(classes)
	count := 1 << uint(n)
	bits := getBits(n)

	ss := make([]string, 1, count)
	ss[0] = ""
	for i := 1; i < count; i++ {
		a := make([]string, 0, n)
		for j := 0; j < n; j++ {
			if (1<<uint(j))&bits[i] != 0 {
				a = append(a, classes[j])
			}
		}
		sort.Strings(a)
		ss = append(ss, "."+strings.Join(a, "."))
	}

	combMemo[key] = ss
	return ss
}

// IntByBits []int sorted by bit count
type IntByBits []int

func (p IntByBits) Len() int { return len(p) }
func (p IntByBits) Less(i, j int) bool {
	if numOfBits(p[i]) == numOfBits(p[j]) {
		return p[i] > p[j]
	}
	return numOfBits(p[i]) < numOfBits(p[j]) //intentional
}
func (p IntByBits) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// count the number of set bits
func numOfBits(bits int) int {
	bits = (bits & 0x55555555) + (bits >> 1 & 0x55555555)
	bits = (bits & 0x33333333) + (bits >> 2 & 0x33333333)
	bits = (bits & 0x0f0f0f0f) + (bits >> 4 & 0x0f0f0f0f)
	bits = (bits & 0x00ff00ff) + (bits >> 8 & 0x00ff00ff)
	return (bits & 0x0000ffff) + (bits >> 16 & 0x0000ffff)
}
