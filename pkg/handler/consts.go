package handler

import (
	"io/ioutil"
	"log"
)

const (
	green = "\x1b[32m"
	cyan  = "\x1b[36m"
	red   = "\x1b[31m"

	end = "\x1b[0m"
)

var logo string

func getLogo() {
	content, err := ioutil.ReadFile("./hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	logo = string(content)
}
