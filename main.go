package main

import (
	"flag"
	"github.com/alexeysoshin/SmsBird/controller"
	"log"
	"net/http"
)

func main() {

	key := flag.String("key", "test_gshuPaZoeEG6ovbc8M79w0QyM", "MessageBird API Key")

	c := controller.NewController(*key)

	m := http.NewServeMux()

	m.HandleFunc("/", c.SendSms)

	log.Println("Server is up and running")
	http.ListenAndServe(":8080", m)
}
