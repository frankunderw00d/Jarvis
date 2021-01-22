package main

import (
	"log"
	"regexp"
)

func init() {}

func main() {
	r := regexp.MustCompile("^[a-zA-X]+[a-zA-Z0-9]{5,17}$")

	log.Println(r.MatchString(""))
}