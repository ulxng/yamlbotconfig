package main

import (
	"fmt"
	"ulxng/yamlbotconf/messages"
)

func main() {
	loader := messages.NewLoader("messages.yaml")
	fmt.Println(loader.GetByKey("psy.faq.intro"))
}
