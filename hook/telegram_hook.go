package hook

import (
	"encoding/json"
	"fmt"
	"github.com/bingoohuang/prettytable"
	"github.com/jxo-me/web-checker/config"
	"github.com/jxo-me/web-checker/core"
	"io"
	"net/http"
	"net/url"
	"time"
)

const TGApi = "https://api.telegram.org/bot%s/sendMessage" // ?chat_id=%d&text=#{ipv4Addr}%0A#{ipv4Result}%0A#{ipv4Domains}

const Tpl = `🚓️服务异常监控预警‼️‼️‼️
<pre>%s</pre>
🔸异常详情: %s
🔸异常时间: %s
🔸监控节点: %s
🔸监控节点IP: %s
`

const NotifyTpl = `🚓️服务监控预警启动
<pre>%s</pre>
🔸监控节点: %s
🔸监控节点IP: %s
🔸采样频率: %s
`

type TelegramHook struct {
	Token   string   `json:"token"`
	ChatId  int64    `json:"chat_id"`
	TimeOut int      `json:"time_out"`
	LocalIp *core.Ip `json:"local_ip"`
}

type Response struct {
	Ok          bool   `json:"ok"`
	Result      Result `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Result struct {
	MessageId int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	Chat      struct {
		Id        int    `json:"id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Entities []struct {
		Offset int    `json:"offset"`
		Length int    `json:"length"`
		Type   string `json:"type"`
	} `json:"entities"`
	From struct {
		Id        int64  `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	} `json:"from"`
}

type Notify struct {
	Name    string `table:"项目名称"`
	Env     string `table:"所属环境"`
	Address string `table:"项目域名"`
	Ssl     string `table:"证书监控"`
}

func (h *TelegramHook) Notify(sites []config.Website, interval int) error {
	var err error
	var list []Notify
	for _, site := range sites {
		list = append(list, Notify{Name: site.Name, Env: site.Env, Address: site.Url, Ssl: "Yes"})
	}
	h.LocalIp, err = core.GetLocation("")
	if err != nil {
		return err
	}
	text := fmt.Sprintf(NotifyTpl, prettytable.TablePrinter{}.Print(&list), h.LocalIp.Country, h.LocalIp.Ip, fmt.Sprintf("%ds", interval))
	// MarkdownV2|HTML|Markdown
	uri := fmt.Sprintf(TGApi, h.Token)
	link := fmt.Sprintf("%s?chat_id=%d&parse_mode=HTML&text=%s", uri, h.ChatId, url.QueryEscape(text))
	return h.SendMsg(link)
}

func (h *TelegramHook) SendMsg(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("cannot new request %s: %w", url, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	if err != nil {
		return fmt.Errorf("cannot set request Header %s: %w", url, err)
	}
	client := http.Client{
		Timeout: time.Second * time.Duration(h.TimeOut),
	}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot request %s: %w", url, err)
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("cannot josn Unmarshal %v: %w", string(body), err)
	}
	if !result.Ok {
		return fmt.Errorf("cannot send telegram message error code: %d, description:%s", result.ErrorCode, result.Description)
	}
	return nil
}

func (h *TelegramHook) Process(resp *core.Response) error {
	if resp.Code != http.StatusOK {
		var content []byte
		if len(resp.Content) > 1 {
			content = resp.Content[1]
		}
		res := core.Result{Name: resp.Website.Name, Env: resp.Website.Env, Address: resp.Website.Url, Status: resp.Code, Elapsed: resp.Duration, Certificate: resp.Certificate}
		text := fmt.Sprintf(Tpl, prettytable.TablePrinter{}.Print(&res), content, time.Now().Format("2006-01-02 15:04:05"), h.LocalIp.Country, h.LocalIp.Ip)
		// MarkdownV2|HTML|Markdown
		uri := fmt.Sprintf(TGApi, h.Token)
		link := fmt.Sprintf("%s?chat_id=%d&parse_mode=HTML&text=%s", uri, h.ChatId, url.QueryEscape(text))
		err := h.SendMsg(link)
		if err != nil {
			return err
		}
	}
	return nil
}
