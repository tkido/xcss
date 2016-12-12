package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var (
	debugFlag bool
	rootFlag  string
	classFlag string

	initClasses []string
)

func init() {
	flag.StringVar(&classFlag, "c", "", "short form of \"class\"")
	flag.StringVar(&classFlag, "class", "", "classes apply to all elements. separator is space e.g. \"foo bar\"")
	flag.BoolVar(&debugFlag, "d", false, "short form of \"debug\"")
	flag.BoolVar(&debugFlag, "debug", false, "add comment where attribute's value came from")
	flag.StringVar(&rootFlag, "r", ".", "short form of \"root\"")
	flag.StringVar(&rootFlag, "root", ".", "path to start recursive walk")
	flag.Parse()

	initClasses = []string{}
	if len(classFlag) > 0 {
		initClasses = strings.Split(classFlag, " ")
	}

	fi, err := os.Stat(rootFlag)
	if err != nil {
		log.Fatal(err)
	}
	if !fi.IsDir() {
		log.Fatal(rootFlag + " is not Directory!!")
	}

}
