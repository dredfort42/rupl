package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
)

var (
	authURL      string
	profileURL   string
	trainingURL  string
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
			loger.Error("Error parsing target URL:", err)
			return
		}
	} else if strings.HasPrefix(r.URL.Path, "/api/v1/profile") {
		targetURL, err = url.Parse("http://" + profileURL)
		if err != nil {
			loger.Error("Error parsing target URL:", err)
			return
		}
	} else if strings.HasPrefix(r.URL.Path, "/api/v1/training") {
		targetURL, err = url.Parse("http://" + trainingURL)
		if err != nil {
			loger.Error("Error parsing target URL:", err)
			return
		}
	} else {
		notFound(w, r)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.Host = targetURL.Host

	// logRequest(r)

	proxy.ServeHTTP(w, r)

	// logResponse(w)
}

// func logRequest(r *http.Request) {
// 	loger.Info("Request Information:", "")
// 	loger.Info("Method:", r.Method)
// 	loger.Info("URL:", r.URL.String())
// 	loger.Info("Proto:", r.Proto)
// 	loger.Info("Host:", r.Host)
// 	loger.Info("Headers:", "")
// 	for key, value := range r.Header {
// 		headerValue := strings.Join(value, ", ")
// 		loger.Info("\t"+key+":", headerValue)
// 	}
// 	loger.Info("------------------------------------------", "")
// }

// func logResponse(w http.ResponseWriter) {
// 	loger.Info("Response Information:", "")
// 	loger.Info("Headers:", "")
// 	for key, value := range w.Header() {
// 		headerValue := strings.Join(value, ", ")
// 		loger.Info("\t"+key+":", headerValue)
// 	}
// 	loger.Info("------------------------------------------", "")
// }

func main() {
	err := cfg.GetConfig()
	if err != nil {
		panic(err)
	}

	authURL = fmt.Sprintf("%s:%s", cfg.Config["auth.host"], cfg.Config["auth.port"])
	profileURL = fmt.Sprintf("%s:%s", cfg.Config["profile.host"], cfg.Config["profile.port"])
	trainingURL = fmt.Sprintf("%s:%s", cfg.Config["training.host"], cfg.Config["training.port"])

	http.HandleFunc("/", sendData)
	http.HandleFunc("/download/", sendData)
	http.HandleFunc("/api/v1/", proxyRequest)

	port := fmt.Sprintf(":%s", cfg.Config["entrypoint.port.ssl"])
	url := fmt.Sprintf("%s://%s%s", cfg.Config["endpoint.protocol.ssl"], cfg.Config["entrypoint.address"], port)
	certFile := "/app/fullchain.pem"
	keyFile := "/app/privkey.pem"

	redirectPortHandler := func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
	}

	// Start HTTP server on port 80
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectPortHandler)); err != nil {
			loger.Error("Error starting HTTP server", err)
			panic(err)
		}
	}()

	// Start HTTPS server on port 443
	if err := http.ListenAndServeTLS(port, certFile, keyFile, nil); err != nil {
		loger.Error("Error starting HTTPS server", err)
		panic(err)
	}

	loger.Success("Entry point", url)
}
