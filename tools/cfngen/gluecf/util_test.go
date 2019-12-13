package gluecf

// utils for testing

import (
	"io/ioutil"
	"os"
)

// Read CF
func readTestFile(filename string) (string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	contents, err := ioutil.ReadAll(fd)
	return string(contents), err
}
