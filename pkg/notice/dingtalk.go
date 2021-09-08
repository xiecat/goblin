package notice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "unknwon.dev/clog/v2"
)

func (d *DingTalk) Send(msg *Msg) error {
	if d.URL == "" {
		return errors.New("DingTalk is not set")
	}
	text := `**访问目标:** %s 上线    
    **访问地址:** %s    
    **地理位置:** %s     
	**规则名称:** %s   
	**访问路径:** %s   
	`
	payload := dingPayload{
		MsgType: "markdown",
		Markdown: dingMarkdown{
			Title: "goblin New Dump",
			Text:  fmt.Sprintf(text, msg.Target, msg.Remote, msg.Location, msg.Rule, msg.Path),
		},
		At: []dingAt{
			{
				IsAtAll: true,
			},
		},
	}
	p, err := json.Marshal(&payload)
	resp, err := http.Post(d.URL, "application/json", strings.NewReader(string(p)))
	if err != nil {
		return fmt.Errorf("HTTP request:%s ", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read HTTP response body: %s", err)
		}
		return fmt.Errorf("non-success response status code %d with body: %s", resp.StatusCode, data)
	}
	return nil
}

func (d *DingTalk) SendTest() {
	if d.URL != "" {
		err := d.Send(&Msg{
			Target:   "this is a test",
			Remote:   "114.114.114.114",
			Location: "中国",
			Rule:     "vpn",
			Path:     "/login.php",
		})
		if err != nil {
			log.Warn("%s", err.Error())
		}
		return
	}
	log.Fatal("DingTalk is Null")
}
