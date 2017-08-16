package main

import (
	"log"

	"github.com/thbourlove/restc"
)

func main() {
	var usernames []string
	restc.NewClient().GetJsonDataWithPath("https://api.github.com/users", &usernames, "$[*].login")
	log.Println(usernames)
}
