package core

import (
	"crypto/tls"
	"github.com/bingoohuang/prettytable"
	"github.com/jxo-me/web-checker/config"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type Processors []func(resp *Response) error

type Checker struct {
	Config     config.Checker
	Processors Processors
	ListResult []Result
	wg         sync.WaitGroup
}

type Result struct {
	Name    string  `table:"项目名称"`
	Env     string  `table:"所属环境"`
	Address string  `table:"系统地址"`
	Status  int     `table:"响应状态"`
	Elapsed float64 `table:"耗时(秒)"`
	Result  string  `table:"请求结果"`
}

func (c *Checker) Run() {
	ticker := time.NewTicker(time.Second * time.Duration(c.Config.Interval))

	for {
		select {
		case <-ticker.C:
			c.ListResult = []Result{}
			for _, site := range c.Config.Websites {
				c.wg.Add(1)
				go func(s config.Website) {
					c.check(s)
					c.wg.Done()
				}(site)
			}
			c.wg.Wait()
			out := prettytable.TablePrinter{}.Print(&c.ListResult)
			log.Printf("\n%s", out)
		}
	}
}

func (c *Checker) check(site config.Website) {
	client := http.Client{
		Timeout: time.Second * time.Duration(c.Config.Timeout),
		Transport: &http.Transport{
			// 注意如果证书已过期，那么只有在关闭证书校验的情况下链接才能建立成功
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := fetchWebsite(&client, site)
	if err != nil {
		log.Printf("unable to fetch the website %s: %v", site.Name, err)
	}
	content := "Success"
	if resp.Code != http.StatusOK {
		if len(resp.Content) > 1 {
			content = string(resp.Content[1])
		}
	}
	res := Result{Name: resp.Website.Name, Env: resp.Website.Env, Address: resp.Website.Url, Status: resp.Code, Elapsed: resp.Duration, Result: content}

	for _, f := range c.Processors {
		err = f(resp)
		if err != nil {
			log.Printf("unable to process response: %v", err)
		}
	}
	c.ListResult = append(c.ListResult, res)
}

func getContent(regex string, body []byte) [][]byte {
	return regexp.MustCompile(regex).FindSubmatch(body)
}
