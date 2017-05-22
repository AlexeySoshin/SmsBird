package main

import (
	"flag"
	"github.com/alexeysoshin/SmsBird/controller"
	"log"
	"net/http"
)

func main() {

	key := flag.String("key", "", "MessageBird API Key")
	flag.Parse()

	c := controller.NewController(*key)

	m := http.NewServeMux()

	m.HandleFunc("/", c.SendSms)

	log.Printf("Server is up and running, key is %s\n", *key)
	http.ListenAndServe(":8080", m)
}
