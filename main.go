package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var help = `Usage: %v FIRST_DIR SECOND_DIR
`

func mapDir(name string) (map[[16]byte]string, error) {
	fullPath, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	dirlist, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	m := make(map[[16]byte]string, len(dirlist))
	for _, item := range dirlist {
		if item.IsDir() {
			continue
		}
		longName := filepath.Join(fullPath, item.Name())
		file, err := os.Open(longName)
		if err != nil {
			return m, err
		}
		data, err := ioutil.ReadAll(file)
		m[md5.Sum(data)] = longName
	}
	return m, nil

}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf(help, os.Args[0])
	}
	local, err := mapDir(os.Args[1])
	remote, err := mapDir(os.Args[2])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	for k := range local {
		if _, ok := remote[k]; ok {
			fmt.Printf("%v\t%v\n", local[k], remote[k])
		}
	}
}
