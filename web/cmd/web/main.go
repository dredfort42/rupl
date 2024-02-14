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

func sendData(w http.ResponseWriter, r *http.Request) {
	path := "/html/index.html"

	if r.URL.Path[1:] != "" {
		path = "/html" + r.URL.Path
	}

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		indexFile, err := os.Open(path + "/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()
		if _, err := io.Copy(w, indexFile); err != nil {
			http.Error(w, "Directory access is forbidden", http.StatusForbidden)
			return
		}
	}

	if strings.Contains(r.URL.Path, "/download/") {
		w.Header().Set("Content-Type", http.DetectContentType(make([]byte, 512))) // Detect content type
		w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	}

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	// log.Printf("[in]\n\t method: %s\n\t scheme %s\n\t host: %s\n\t path: %s\n\t body: %s\n\t requestURI: %s\n\t rawQuery: %s\n\t fragment: %s\n\t opaque: %s\n\t user: %s\n\t rawPath: %s\n\t forceQuery: %t\n\t rawFragment: %s\n",
	// 	r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, r.Body, r.RequestURI, r.URL.RawQuery, r.URL.Fragment, r.URL.Opaque, r.URL.User, r.URL.RawPath, r.URL.ForceQuery, r.URL.RawFragment)

	request, err := http.NewRequest(r.Method, r.URL.Path, nil)
	if err != nil {
		logprinter.PrintError("Error creating request", err)
		return
	}

	request.URL = r.URL
	request.Body = r.Body
	request.Header = r.Header
	request.URL.Scheme = "http"

	if strings.HasPrefix(r.URL.Path, "/api/v1/auth/") {
		request.URL.Host = authURL
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		logprinter.PrintError("Error sending request: %v\n", err)
		return
	}
	defer response.Body.Close()

	logprinter.PrintInfo(response.Status, request.URL.Host+request.URL.Path+request.URL.RawQuery)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logprinter.PrintError("Error reading response body: %v\n", err)
		return
	}

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

	http.HandleFunc("/", sendData)
	http.HandleFunc("/api/v1/auth/", proxyRequest)

	port := fmt.Sprintf(":%s", config["entrypoint.port.ssl"])
	url := fmt.Sprintf("%s://%s%s", config["entrypoint.protocol.ssl"], config["entrypoint.address"], port)
	certFile := "/app/fullchain.pem"
	keyFile := "/app/privkey.pem"

	logprinter.PrintSuccess("Entry point", url)
	logprinter.PrintError("Connection error:", http.ListenAndServeTLS(port, certFile, keyFile, nil))
}
