package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var (
	debugFlag  bool
	watchFlag  bool
	deleteFlag bool
	rootFlag   string
	classFlag  string

	initClasses []string
)

func init() {
	flag.BoolVar(&debugFlag, "d", false, "short form of \"debug\"")
	flag.BoolVar(&debugFlag, "debug", false, "add comment where attribute's value came from")
	flag.BoolVar(&watchFlag, "w", false, "short form of \"watch\"")
	flag.BoolVar(&watchFlag, "watch", false, "watch \"xcss\" and \"sxml\" files and run convert when these files are changed")
	flag.StringVar(&classFlag, "c", "", "short form of \"class\"")
	flag.StringVar(&classFlag, "class", "", "classes apply to all elements. separator is space e.g. \"foo bar\"")
	flag.StringVar(&rootFlag, "r", ".", "short form of \"root\"")
	flag.StringVar(&rootFlag, "root", ".", "path to start recursive walk")
	flag.StringVar(&deleteFlag, "delete", false, "delete all XCSS and SXML files after conversion")
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
