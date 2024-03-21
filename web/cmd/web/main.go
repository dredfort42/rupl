package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
)

var (
	authURL      string
	profileURL   string
	notFoundPath string = "/html/404.html"
)

func notFound(w http.ResponseWriter, _ *http.Request) {
	notFoundFile, err := os.Open(notFoundPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer notFoundFile.Close()

	if _, err := io.Copy(w, notFoundFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendData(w http.ResponseWriter, r *http.Request) {
	path := "/html/index.html"

	if r.URL.Path[1:] != "" {
		path = "/html" + r.URL.Path
	}

	file, err := os.Open(path)
	if err != nil {
		notFound(w, r)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		notFound(w, r)
		return
	}

	if fileInfo.IsDir() {
		indexFile, err := os.Open(path + "/index.html")
		if err != nil {
			notFound(w, r)
			return
		}
		defer indexFile.Close()

		if _, err := io.Copy(w, indexFile); err != nil {
			notFound(w, r)
			return
		}

		return
	}

	if strings.Contains(r.URL.Path, "/download/") {
		w.Header().Set("Content-Type", http.DetectContentType(make([]byte, 512))) // Detect content type
		w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	}

	if _, err = io.Copy(w, file); err != nil {
		notFound(w, r)
		return
	}
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	var targetURL *url.URL
	var err error

	if strings.HasPrefix(r.URL.Path, "/api/v1/auth/") {
		targetURL, err = url.Parse("http://" + authURL)
		if err != nil {
			logprinter.PrintError("Error parsing target URL:", err)
			return
		}
	} else if strings.HasPrefix(r.URL.Path, "/api/v1/profile") {
		targetURL, err = url.Parse("http://" + profileURL)
		if err != nil {
			logprinter.PrintError("Error parsing target URL:", err)
			return
		}
	} else {
		notFound(w, r)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.Host = targetURL.Host

	if DEBUG {
		logRequest(r)
	}

	proxy.ServeHTTP(w, r)

	if DEBUG {
		logResponse(w)
	}

}

func logRequest(r *http.Request) {
	logprinter.PrintInfo("Request Information:", "")
	logprinter.PrintInfo("Method:", r.Method)
	logprinter.PrintInfo("URL:", r.URL.String())
	logprinter.PrintInfo("Proto:", r.Proto)
	logprinter.PrintInfo("Host:", r.Host)
	logprinter.PrintInfo("Headers:", "")
	for key, value := range r.Header {
		headerValue := strings.Join(value, ", ")
		logprinter.PrintInfo("\t"+key+":", headerValue)
	}
	logprinter.PrintInfo("------------------------------------------", "")
}

func logResponse(w http.ResponseWriter) {
	logprinter.PrintInfo("Response Information:", "")
	logprinter.PrintInfo("Headers:", "")
	for key, value := range w.Header() {
		headerValue := strings.Join(value, ", ")
		logprinter.PrintInfo("\t"+key+":", headerValue)
	}
	logprinter.PrintInfo("------------------------------------------", "")
}

var DEBUG bool = false

func main() {

	if debug := os.Getenv("DEBUG"); debug != "" {
		logprinter.PrintWarning("Debugging is enabled.", "")
		DEBUG = true
	}

	config, err := configreader.GetConfig()

	if err != nil {
		os.Exit(1)
	}

	authURL = fmt.Sprintf("%s:%s", config["auth.host"], config["auth.port"])
	profileURL = fmt.Sprintf("%s:%s", config["profile.host"], config["profile.port"])

	http.HandleFunc("/", sendData)
	http.HandleFunc("/download/", sendData)
	http.HandleFunc("/api/v1/", proxyRequest)

	port := fmt.Sprintf(":%s", config["entrypoint.port.ssl"])
	url := fmt.Sprintf("%s://%s%s", config["entrypoint.protocol.ssl"], config["entrypoint.address"], port)
	certFile := "/app/fullchain.pem"
	keyFile := "/app/privkey.pem"

	redirectPortHandler := func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
	}

	// Start HTTP server on port 80
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectPortHandler)); err != nil {
			logprinter.PrintError("Error starting HTTP server", err)
			panic(err)
		}
	}()

	// Start HTTPS server on port 443
	if err := http.ListenAndServeTLS(port, certFile, keyFile, nil); err != nil {
		logprinter.PrintError("Error starting HTTPS server", err)
		panic(err)
	}

	logprinter.PrintSuccess("Entry point", url)
}
