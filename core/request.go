package core

import (
	"fmt"
	"github.com/jxo-me/web-checker/config"
	"io"
	"log"
	"net/http"
	"time"
)

func fetchWebsite(client *http.Client, site config.Website) (*Response, error) {
	// for more granular metrics you may use https://golang.org/pkg/net/http/httptrace/
	start := time.Now()
	req, err := http.NewRequest("GET", site.Url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	if err != nil {
		return nil, fmt.Errorf("cannot http NewRequest %s: %w", site.Url, err)
	}
	resp, err := client.Do(req)
	elapsed := time.Since(start).Seconds()

	if err != nil {
		return nil, fmt.Errorf("cannot fetch website %s: %w", site.Name, err)
	}
	defer func() { _ = resp.Body.Close() }()
	statusCode := resp.StatusCode
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var content [][]byte
	if site.Regex != "" {
		content = getContent(site.Regex, bytes)
	}
	certificate := "正常"
	// check ssl certificate expired
	for _, cert := range resp.TLS.PeerCertificates {
		if !cert.NotAfter.After(time.Now()) {
			statusCode = http.StatusBadRequest
			msg := fmt.Sprintf("Website [%s] certificate has expired: %s", site.Url, cert.NotAfter.Local().Format("2006-01-02 15:04:05"))
			log.Println(msg)
			content = [][]byte{}
			content = append(content, []byte("Expired"))
			content = append(content, []byte(fmt.Sprintf("ssl certificate has expired: %s", cert.NotAfter.Local().Format("2006-01-02 15:04:05"))))
			certificate = "过期"
		}
	}

	return &Response{
		Website:     site,
		Code:        statusCode,
		Duration:    elapsed,
		Content:     content,
		Certificate: certificate,
	}, nil
}
