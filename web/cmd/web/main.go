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
	// log.Printf("[in]\n\t method: %s\n\t scheme %s\n\t host: %s\n\t path: %s\n\t body: %s\n\t requestURI: %s\n\t rawQuery: %s\n\t fragment: %s\n\t opaque: %s\n\t user: %s\n\t rawPath: %s\n\t forceQuery: %t\n\t rawFragment: %s\n",
	// 	r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, r.Body, r.RequestURI, r.URL.RawQuery, r.URL.Fragment, r.URL.Opaque, r.URL.User, r.URL.RawPath, r.URL.ForceQuery, r.URL.RawFragment)

	log.Printf("[in]\n\t method: %s\n\t scheme %s\n\t host: %s\n\t path: %s\n\t requestURI: %s\n\t forceQuery: %t\n",
		r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, r.RequestURI, r.URL.RawQuery, r.URL.ForceQuery)

	if r.URL.Path[1:] == "" {
		fmt.Fprintln(w, readFile("./html/index.html"))
	} else {
		filePath := "./html" + r.URL.Path

		if _, err := os.Stat(filePath); err == nil {
			fmt.Fprintln(w, readFile(filePath))
		} else if os.IsNotExist(err) {
			fmt.Fprintln(w, readFile("./html/index.html"))
		} else {
			logprinter.PrintError("Error checking file existence: %v\n", err)
		}

	}
}

func deviceAuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[in]\n\t method: %s\n\t scheme %s\n\t host: %s\n\t path: %s\n\t requestURI: %s\n\t forceQuery: %t\n",
		r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, r.RequestURI, r.URL.RawQuery, r.URL.ForceQuery)

	respond := `
	{
    	"device_code": "GmRhmhcxhwAzkoEqiMEg_DnyEysNkuNhszIySk9eS",
        "user_code": "WDJB-MJHT",
        "verification_uri": "https://example.com/device",
        "verification_uri_complete":
            "https://example.com/device?user_code=WDJB-MJHT",
        "expires_in": 1800,
        "interval": 5
    }`

	fmt.Fprintln(w, respond)
}

func main() {
	config, err := configreader.GetConfig()

	if err != nil {
		os.Exit(1)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/device_authorization", deviceAuthHandler)

	// port := fmt.Sprintf(":%s", config["entrypoint.port"])
	// url := fmt.Sprintf("%s://%s%s", config["entrypoint.protocol"], config["entrypoint.address"], port)

	// logprinter.PrintSuccess("Entry point", url)
	// log.Fatal(http.ListenAndServe(port, nil))

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
