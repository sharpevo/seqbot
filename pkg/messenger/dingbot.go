package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type DingBot struct {
	Token string
}

func NewDingBot(token string) *DingBot {
	return &DingBot{
		Token: token,
	}
}

func (d *DingBot) String() string {
	return fmt.Sprintf("DingBot %s", d.Token)
}

func (d *DingBot) Url() string {
	return fmt.Sprintf(
		"https://oapi.dingtalk.com/robot/send?access_token=%s", d.Token)
}

func (d *DingBot) Send(message string) error {
	data := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "seqbot",
			"text":  message,
		},
	}
	dataBytes, _ := json.Marshal(data)
	logrus.Debugf("send message to dingtalk: %s", string(dataBytes))
	req, _ := http.NewRequest("POST", d.Url(), bytes.NewBuffer(dataBytes))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugf("response from dingtalk: %s", string(body))
	return nil
}
