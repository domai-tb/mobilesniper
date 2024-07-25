package enum

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func IsHTTPorHTTPS(ip string, port int) bool {
	urlHTTP := fmt.Sprintf("http://%s:%d", ip, port)
	urlHTTPS := fmt.Sprintf("https://%s:%d", ip, port)

	// Check HTTP
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(urlHTTP)
	if err == nil {
		resp.Body.Close()
		return true
	}

	// Check HTTPS
	client = &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err = client.Get(urlHTTPS)
	if err == nil {
		resp.Body.Close()
		return true
	}

	return false
}
