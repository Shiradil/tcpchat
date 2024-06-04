package handler

import (
	"io/ioutil"
	"log"
)

var logo string

func getLogo() {
	content, err := ioutil.ReadFile("./hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	logo = string(content)
}
