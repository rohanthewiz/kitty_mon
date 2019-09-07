package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var apiKey, apiSecret, smsTo, smsFrom string

func init() {
	apiKey = os.Getenv("SMS_KEY")
	apiSecret = os.Getenv("SMS_SECRET")
	smsTo = os.Getenv("SMS_TO")
	smsFrom = os.Getenv("SMS_FROM")
}

type NexmoSmsResp struct {
	MessageCount string `json:"message_count"`
	Messages     []struct {
		To               string `json:"to"`
		MessageId        string `json:"message-id"`
		Status           string `json:"status"`
		RemainingBalance string `json:"remaining-balance"`
		MessagePrice     string `json:"message-price"`
		Network          string `json:"network"`
	} `json:"messages"`
}

const nexmoSendURL = "https://rest.nexmo.com/sms/json"

func NexmoSend(msg string) (err error) {
	//errStage := " when sending nexmo message"
	if apiKey == "" || apiSecret == "" || smsTo == "" || smsFrom == "" {
		return errors.New("empty SMS_KEY or SMS_SECRET or SMS_TO or SMS_FROM")
	}

	fd := url.Values{}
	fd.Set("to", smsTo)
	fd.Set("from", smsFrom)
	fd.Set("text", msg)
	fd.Set("api_key", apiKey)
	fd.Set("api_secret", apiSecret)
	fdReader := strings.NewReader(fd.Encode())

	req, _ := http.NewRequest("POST", nexmoSendURL, fdReader)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err // todo ! - Wrap
	}
	if resp.StatusCode >= 400 { // todo - log
		// logger.LogErr(serr.New("Status code not ok" + errStage, "status_code", strconv.Itoa(resp.StatusCode),
		// 	"endpoint", publicKeysEndpoint))
		return errors.New("Status code not ok" + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()
	strBody, err := ioutil.ReadAll(resp.Body)
	rsp := NexmoSmsResp{}
	err = json.Unmarshal(strBody, &rsp)
	if err != nil {
		// logger.LogErr(err, "Unable to parse JSON response" + errStage) // todo - log
		return err
	}

	fmt.Printf("rsp %#v\n", rsp)

	return
}
