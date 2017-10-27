package mail

import (
	"bytes"
	"datamesh.com/common/utils/web"
	"encoding/json"
	"fmt"
)

func Send(address, from, fromName, to, subject, html string) error {
	scm := sendCloudMail{
		From:     from,
		Fromname: fromName,
		To:       to,
		Subject:  subject,
		Html:     html,
	}
	return scm.send(address)
}

// the struct holds the neccessary info needed by SendCloud Email ISP.
type sendCloudMail struct {
	From     string `json:"from"`
	Fromname string `json:"fromname"`
	To       string `json:"to"` // we do not support multiple reciever!
	Subject  string `json:"subject"`
	Html     string `json:"html"`
}

type emailResp struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type EmailSuccessJson struct {
	Checksum    string `json:"checksum"`
	Fingerprint string `json:"fingerprint"`
}

func (mail *sendCloudMail) send(address string) error {
	f := func(msg string) error {
		return fmt.Errorf("Alert mail: send failed, message: \n%v", msg)
	}
	b, err := json.Marshal(mail)
	if err != nil {
		return f(err.Error())
	}
	wp := web.WebPage{
		Url:      address,
		Method:   "POST",
		BodyType: "application/json",
		Body:     bytes.NewBuffer(b),
	}
	if err := wp.DoRequest(); err != nil {
		return f(err.Error())
	}
	if wp.StatusCode != 200 {
		return f(string(wp.RespBody))
	}
	r := emailResp{}
	if err := json.Unmarshal(wp.RespBody, &r); err != nil {
		return f(err.Error())
	}
	if r.Code != 1 && r.Code != 10200 { // if the mail service is refactored, must modify this.
		return f(r.Message)
	}
	return nil
}
