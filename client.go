package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TextClient interface {
	SendMessage(string) error
}

type StubClient struct{}

func (c *StubClient) SendMessage(msg string) error {
	fmt.Printf("Sending message '%s'\n", msg)
	return nil
}

type TwilioConfig struct {
	AccountSid string
	AuthToken  string
	UrlStr     string
	BfNumber   string
	FromNumber string
}

func NewTwilioClient(AccountSid string, AuthToken string, BfNumber string, FromNumber string) *TwilioClient {
	UrlStr := "https://api.twilio.com/2010-04-01/Accounts/" + AccountSid + "/Messages.json"
	config := &TwilioConfig{
		AccountSid: AccountSid,
		AuthToken:  AuthToken,
		UrlStr:     UrlStr,
		BfNumber:   BfNumber,
		FromNumber: FromNumber,
	}

	return &TwilioClient{Config: config}
}

type TwilioClient struct {
	Config *TwilioConfig
}

func (tc *TwilioClient) SendMessage(msg string) error {
	msgData := url.Values{}
	msgData.Set("To", tc.Config.BfNumber)
	msgData.Set("From", tc.Config.FromNumber)
	msgData.Set("Body", msg)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", tc.Config.UrlStr, &msgDataReader)
	if err != nil {
		fmt.Println("Error creating new request to twilio")
		return err
	}
	req.SetBasicAuth(tc.Config.AccountSid, tc.Config.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to twilio")
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
	return nil
}
