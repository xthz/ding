package ding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Webhook struct {
	AccessToken string
	Secret      string
}

type response struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

//SendMessageText Function to send message
//goland:noinspection GoUnhandledErrorResult
func (t *Webhook) SendMessageText(text string, at ...string) error {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": text,
		},
	}

	if len(at) == 0 {
	} else if len(at) == 1 {
		if at[0] == "*" { // at all
			msg["at"] = map[string]interface{}{
				"isAtAll": true,
			}
		} else { // at specific user
			re := regexp.MustCompile(`^\+*\d{10,15}$`)
			if re.MatchString(at[0]) {
				msg["at"] = map[string]interface{}{
					"atMobiles": at,
					"isAtAll":   false,
				}
			} else {
				return errors.New(`parameter error, "at" parameter must be in "*" or mobile phone number format`)
			}
		}
	} else {
		re := regexp.MustCompile(`^\+*\d{10,15}$`)
		for _, v := range at {
			if !re.MatchString(v) {
				return errors.New(`parameter error, "at" parameter must be in "*" or mobile phone number format`)
			}
		}
		msg["at"] = map[string]interface{}{
			"atMobiles": at,
			"isAtAll":   false,
		}
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(t.getURL(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New(fmt.Sprintf("response error: %s", string(body)))
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (t *Webhook) sendMessageMarkdown(title, text string, at ...string) error {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
	}

	if len(at) == 0 {
	} else if len(at) == 1 {
		if at[0] == "*" { // at all
			msg["at"] = map[string]interface{}{
				"isAtAll": true,
			}
		} else { // at specific user
			re := regexp.MustCompile(`^\+*\d{10,15}$`)
			if re.MatchString(at[0]) {
				msg["at"] = map[string]interface{}{
					"atMobiles": at,
					"isAtAll":   false,
				}
			} else {
				return errors.New(`parameter error, "at" parameter must be in "*" or mobile phone number format`)
			}
		}
	} else {
		re := regexp.MustCompile(`^\+*\d{10,15}$`)
		for _, v := range at {
			if !re.MatchString(v) {
				return errors.New(`parameter error, "at" parameter must be in "*" or mobile phone number format`)
			}
		}
		msg["at"] = map[string]interface{}{
			"atMobiles": at,
			"isAtAll":   false,
		}
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(t.getURL(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func (t *Webhook) hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (t *Webhook) getURL() string {
	wh := "https://oapi.dingtalk.com/robot/send?access_token=" + t.AccessToken
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, t.Secret)
	sign := t.hmacSha256(stringToSign, t.Secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", wh, timestamp, sign)
	return url
}
