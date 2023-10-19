package core

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Ip struct {
	Ip          string  `json:"ip"`
	Asn         string  `json:"asn"`
	Netmask     int     `json:"netmask"`
	Hostname    string  `json:"hostname"`
	City        string  `json:"city"`
	PostCode    string  `json:"post_code"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type IpResp struct {
	Ip Ip `json:"ip"`
}

// GetLocation
// https://ip.nf/me.json
// https://ip.nf/208.80.154.224.json
// https://ip.seeip.org/geoip
// https://ip.seeip.org/geoip/1.1.1.1
func GetLocation(ip string) (*Ip, error) {
	url := "https://ip.nf/me.json"
	if ip != "" {
		url = fmt.Sprintf("https://ip.nf/%s.json", ip)
	}
	client := http.Client{
		Transport: &http.Transport{
			// 注意如果证书已过期，那么只有在关闭证书校验的情况下链接才能建立成功
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	var result IpResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("cannot josn Unmarshal %v: %w", string(body), err)
	}

	return &result.Ip, nil
}
