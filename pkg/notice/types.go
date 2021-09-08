package notice

type DingTalk struct {
	URL string `yaml:"DingTalk"`
}

type dingAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type dingMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type dingPayload struct {
	MsgType  string       `json:"msgtype"`
	Markdown dingMarkdown `json:"markdown"`
	At       []dingAt     `json:"at"`
}

type Msg struct {
	Target      string
	Rule        string
	Path        string
	Remote      string
	Location    string
	RequestDump string
}
