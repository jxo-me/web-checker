package core

import (
	"encoding/json"
	"fmt"
	"github.com/jxo-me/web-checker/config"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func fetchWebsite(client *http.Client, site config.Website) (*Response, error) {
	// for more granular metrics you may use https://golang.org/pkg/net/http/httptrace/
	start := time.Now()
	postPara := ""
	method := "GET"
	contentType := "application/x-www-form-urlencoded"
	if site.Method != "" {
		method = strings.ToUpper(site.Method)
	}
	if site.Body != "" {
		method = "POST"
		postPara = site.Body
		if json.Valid([]byte(postPara)) {
			contentType = "application/json"
			// 如果 RequestBody 的 JSON 无效但前缀为 JSON 括号则为 JSON
		} else if hasJSONPrefix(postPara) {
			log.Println("RequestBody 的 JSON 无效！")
			return nil, fmt.Errorf("the JSON for RequestBody is invalid %s", postPara)
		}
	}
	req, err := http.NewRequest(method, site.Url, strings.NewReader(postPara))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	if err != nil {
		return nil, fmt.Errorf("cannot http NewRequest %s: %w", site.Url, err)
	}
	headers := CheckParseHeaders(site.Headers)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	req.Header.Add("content-type", contentType)

	resp, err := client.Do(req)
	elapsed := time.Since(start).Seconds()

	var content [][]byte
	statusCode := resp.StatusCode
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
	if err != nil {
		content = append(content, []byte("Error"))
		content = append(content, []byte(err.Error()))
		return &Response{
			Website:     site,
			Code:        http.StatusGatewayTimeout,
			Duration:    elapsed,
			Content:     content,
			Certificate: certificate,
		}, fmt.Errorf("cannot fetch website %s: %w", site.Name, err)
	}
	defer func() { _ = resp.Body.Close() }()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if site.Regex != "" {
		content = getContent(site.Regex, bytes)
	}

	return &Response{
		Website:     site,
		Code:        statusCode,
		Duration:    elapsed,
		Content:     content,
		Certificate: certificate,
	}, nil
}

func CheckParseHeaders(headerStr string) (headers map[string]string) {
	headers = make(map[string]string)
	headerArr := strings.Split(headerStr, "\r\n")
	for _, header := range headerArr {
		header = strings.TrimSpace(header)
		if header != "" {
			parts := strings.Split(header, ":")
			if len(parts) != 2 {
				log.Printf("Header不正确: %s", header)
				continue
			}
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return headers
}

// hasJSONPrefix returns true if the string starts with a JSON open brace.
func hasJSONPrefix(s string) bool {
	return strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[")
}
