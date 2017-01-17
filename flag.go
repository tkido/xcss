package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

// Flags keeps settings from commandline
type Flags struct {
	Debug   bool
	Delete  bool
	Root    string
	Watch   bool
	Classes []string
}

var flags Flags

func init() {
	var class string

	flag.StringVar(&class, "c", "", "short form of \"class\"")
	flag.StringVar(&class, "class", "", "classes apply to all elements. separator is space e.g. \"foo bar\"")
	flag.BoolVar(&flags.Debug, "d", false, "short form of \"debug\"")
	flag.BoolVar(&flags.Debug, "debug", false, "add comment where attribute's value came from")
	flag.BoolVar(&flags.Delete, "delete", false, "delete all XCSS and SXML files after conversion")
	flag.StringVar(&flags.Root, "r", ".", "short form of \"root\"")
	flag.StringVar(&flags.Root, "root", ".", "path to start recursive walk")
	flag.BoolVar(&flags.Watch, "w", false, "short form of \"watch\"")
	flag.BoolVar(&flags.Watch, "watch", false, "watch \"xcss\" and \"sxml\" files and run convert when these files are changed")
	flag.Parse()

	if len(class) > 0 {
		flags.Classes = strings.Split(class, " ")
	}

	fi, err := os.Stat(flags.Root)
	if err != nil {
		log.Fatal(err)
	}
	if !fi.IsDir() {
		log.Fatal(flags.Root + " is not Directory!!")
	}

}
