package client

import (
	ku "bitpay/key_utils"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

type Client struct {
	ApiUri   string
	Insecure bool
	ClientId string
}

func (client *Client) CreateInvoice(price float64, currency string) (invoice map[string]string, err error) {
	match, _ := regexp.MatchString("^[[:upper:]]{3}$", currency)
	if !match {
		err = errors.New("BitPayArgumentError: invalid currency code")
		return nil, err
	}
	url := client.ApiUri + "/invoices"
	htclient := setHttpClient(client)
	paylo := make(map[string]string)
	payload, _ := json.Marshal(paylo)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	response, _ := htclient.Do(req)
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	var jsonContents map[string]interface{}
	json.Unmarshal(contents, &jsonContents)
	if response.StatusCode/100 != 2 {
		err = processErrorMessage(response, jsonContents)
	} else {
		invoice, err = processInvoice(response, jsonContents)
		err = nil
	}
	return invoice, err
}

func (client *Client) PairWithCode(str string) (token map[string]string, err error) {
	pm := ku.GeneratePem()
	sin := ku.GenerateSinFromPem(pm)
	client.ClientId = sin
	token = make(map[string]string)
	match, _ := regexp.MatchString("^[[:alnum:]]{7}$", str)
	if !match {
		token = nil
		err = errors.New("BitPayArgumentError: invalid pairing code")
		return token, err
	}
	url := client.ApiUri + "/tokens"
	htclient := setHttpClient(client)
	paylo := make(map[string]string)
	paylo["id"] = client.ClientId
	paylo["pairingCode"] = str
	payload, _ := json.Marshal(paylo)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	response, _ := htclient.Do(req)
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	var jsonContents map[string]interface{}
	json.Unmarshal(contents, &jsonContents)
	if response.StatusCode/100 != 2 {
		err = processErrorMessage(response, jsonContents)
	} else {
		token, err = processToken(response, jsonContents)
		err = nil
	}
	return token, err
}

func setHttpClient(client *Client) *http.Client {
	if client.Insecure {
		trans := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		htclient := &http.Client{Transport: trans}
		return htclient
	} else {
		trans := &http.Transport{}
		htclient := &http.Client{Transport: trans}
		return htclient
	}
}

func processErrorMessage(response *http.Response, jsonContents map[string]interface{}) error {
	responseStatus := strconv.Itoa(response.StatusCode)
	contentError := responseStatus + ": " + jsonContents["error"].(string)
	return errors.New(contentError)
}

func processInvoice(response *http.Response, jsonContents map[string]interface{}) (token map[string]string, err error) {
	datarray := jsonContents["data"].(interface{})
	data := datarray.(map[string]interface{})
	token, _ = processSuccess(data)
	return token, nil
}

func processToken(response *http.Response, jsonContents map[string]interface{}) (invoice map[string]string, err error) {
	datarray := jsonContents["data"].([]interface{})
	data := datarray[0].(map[string]interface{})
	invoice, _ = processSuccess(data)
	return invoice, nil
}

func processSuccess(data map[string]interface{}) (token map[string]string, err error) {
	token = make(map[string]string)
	for k := range data {
		if str, ok := data[k].(string); ok {
			token[k] = str
		} else {
			notString := data[k]
			token[k] = fmt.Sprintf("%s", notString)
		}
	}
	return token, nil
}
