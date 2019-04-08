package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/singmyr/astrocyte-bot/slack"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Error parsing the body")
			return
		}

		rData, _ := slack.DataFromBytes(b)

		slack.Handle(w, rData)
	}
}

func main() {
	slack.RegisterCommand(&slack.Command{
		Command: "ngrok",
		Handler: func(w io.Writer, d *slack.RequestData) {
			fmt.Fprintf(w, "Wooh, inside ngrok! - %s c: %s", d.UserName, d.Command)
		},
	})
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":4390", nil))
}
