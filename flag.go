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
	flag.BoolVar(&debugFlag, "debug", false, "debug flag")
	flag.BoolVar(&debugFlag, "d", false, "debug flag")
	flag.StringVar(&rootFlag, "root", ".", "root flag")
	flag.StringVar(&rootFlag, "r", ".", "root flag")
	flag.StringVar(&classFlag, "class", "", "class flag")
	flag.StringVar(&classFlag, "c", "", "class flag")
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
