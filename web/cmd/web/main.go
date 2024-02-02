package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
)

var (
	authURL string
)

func readFile(path string) string {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(fileContent)
}

func showContent(w http.ResponseWriter, r *http.Request) {
	// log.Printf("[in]\n\t method: %s\n\t scheme %s\n\t host: %s\n\t path: %s\n\t body: %s\n\t requestURI: %s\n\t rawQuery: %s\n\t fragment: %s\n\t opaque: %s\n\t user: %s\n\t rawPath: %s\n\t forceQuery: %t\n\t rawFragment: %s\n",
	// 	r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, r.Body, r.RequestURI, r.URL.RawQuery, r.URL.Fragment, r.URL.Opaque, r.URL.User, r.URL.RawPath, r.URL.ForceQuery, r.URL.RawFragment)

	// log.Printf("[in]\n\t method: %s\n\t scheme %s\n\t host: %s\n\t path: %s\n\t requestURI: %s\n\t forceQuery: %t\n",
	// 	r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, r.RequestURI, r.URL.RawQuery, r.URL.ForceQuery)

	if r.URL.Path[1:] == "" {
		fmt.Fprintln(w, readFile("./html/index.html"))
	} else {
		path := "./html" + r.URL.Path

		if _, err := os.Stat(path); err == nil {
			fmt.Fprintln(w, readFile(path))
		} else if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, readFile("./html/index.html"))
		} else {
			logprinter.PrintError("Error checking file existence: %v\n", err)
		}

	}
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	//log.Printf("[in]\n\t method: %s\n\t scheme %s\n\t host: %s\n\t path: %s\n\t body: %s\n\t requestURI: %s\n\t rawQuery: %s\n\t fragment: %s\n\t opaque: %s\n\t user: %s\n\t rawPath: %s\n\t forceQuery: %t\n\t rawFragment: %s\n",
	//	r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, r.Body, r.RequestURI, r.URL.RawQuery, r.URL.Fragment, r.URL.Opaque, r.URL.User, r.URL.RawPath, r.URL.ForceQuery, r.URL.RawFragment)

	request, err := http.NewRequest(r.Method, r.URL.Path, nil)
	if err != nil {
		logprinter.PrintError("Error creating request", err)
		return
	}

	request.URL.Scheme = "http"
	if strings.HasPrefix(r.URL.Path, "/api/v1/auth/") {
		request.URL.Host = authURL
	}

	request.Body = r.Body

	logprinter.PrintInfo("Request URL", request.URL.Host+request.URL.Path)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		logprinter.PrintError("Error sending request: %v\n", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logprinter.PrintError("Error reading response body: %v\n", err)
		return
	}

	logprinter.PrintInfo("Response status", response.Status)
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))

	if _, err := w.Write(body); err != nil {
		panic(err)
	}
}

func main() {
	config, err := configreader.GetConfig()

	if err != nil {
		os.Exit(1)
	}

	authURL = fmt.Sprintf("%s:%s", config["auth.host"], config["auth.port"])

	http.HandleFunc("/", showContent)
	http.HandleFunc("/api/v1/auth/", proxyRequest)

	// Successfully received certificate.
	// Certificate is saved at: /etc/letsencrypt/live/rupl.org/fullchain.pem
	// Key is saved at:         /etc/letsencrypt/live/rupl.org/privkey.pem
	// This certificate expires on 2024-04-29.
	// These files will be updated when the certificate renews.
	// Certbot has set up a scheduled task to automatically renew this certificate in the background.

	port := fmt.Sprintf(":%s", config["entrypoint.port.ssl"])
	url := fmt.Sprintf("%s://%s%s", config["entrypoint.protocol.ssl"], config["entrypoint.address"], port)
	certFile := "./ssl/fullchain.pem"
	keyFile := "./ssl/privkey.pem"

	logprinter.PrintSuccess("Entry point", url)
	log.Fatal(http.ListenAndServeTLS(port, certFile, keyFile, nil))
}
