package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var inputFile, outputFile, packageName, root string
	flag.StringVar(&inputFile, "in", "tags.json", "input file name")
	flag.StringVar(&outputFile, "out", "tags.go", "output file name")
	flag.StringVar(&packageName, "p", "main", "Go package name")
	flag.StringVar(&root, "root", "", "Root node to start walking from")
	flag.Parse()

	dat, err := ioutil.ReadFile(inputFile)
	check(err)

	structs, err := transform(dat, root)
	check(err)

	fileContents := addFileMeta(packageName, structs)

	fp := filepath.Dir(outputFile)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		os.MkdirAll(fp, os.ModePerm)
	}

	err = ioutil.WriteFile(outputFile, []byte(fileContents), 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func addFileMeta(packageName, structs string) string {
	return fmt.Sprintf("package %s\n\nconst (\n%s\n)", packageName, structs)
}
