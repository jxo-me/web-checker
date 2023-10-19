package core

import (
	"encoding/json"
	"github.com/jxo-me/web-checker/config"
)

type Response struct {
	Website  config.Website
	Code     int      `json:"code"`
	Duration float64  `json:"duration"`
	Content  [][]byte `json:"content"`
}

func (r *Response) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func FromJson(data []byte) (*Response, error) {
	resp := new(Response)
	err := json.Unmarshal(data, resp)
	return resp, err
}
