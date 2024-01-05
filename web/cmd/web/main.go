package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
)

func readFile(path string) string {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(fileContent)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" {
		fmt.Fprintln(w, readFile("./html/index.html"))
	} else {
		fmt.Fprintln(w, readFile("./html"+r.URL.Path))
	}
}

func main() {
	config, err := configreader.GetConfig()

	if err != nil {
		os.Exit(1)
	}

	http.HandleFunc("/", homeHandler)

	port := fmt.Sprintf(":%s", config["entrypoint.port"])
	url := fmt.Sprintf("%s://%s%s", config["entrypoint.protocol"], config["entrypoint.address"], port)

	logprinter.PrintSuccess("Entry point", url)
	log.Fatal(http.ListenAndServe(port, nil))
}
