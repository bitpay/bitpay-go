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
	"strings"
)

type Client struct {
	ApiUri   string
	Insecure bool
	ClientId string
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
	if client.Insecure {
		url := client.ApiUri + "/tokens"
		trans := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		htclient := &http.Client{Transport: trans}
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
		if response.StatusCode/100 != 2 {
			fmt.Println("we should be here")
			json.Unmarshal(contents, &jsonContents)
			responseStatus := strconv.Itoa(response.StatusCode)
			contentError := responseStatus + fmt.Sprintf(": %s", jsonContents["error"])
			err = errors.New(contentError)
			token["facade"] = contentError
		} else {
			fmt.Println("what are we doing here")
			json.Unmarshal(contents, &jsonContents)
			datarray := jsonContents["data"].([]interface{})
			data := datarray[0].(map[string]interface{})
			token["facade"] = fmt.Sprintf("%s", data["facade"])
			token["token"] = fmt.Sprintf("%s", data["token"])
			err = nil
		}
	} else {
		url := client.ApiUri + "/tokens"
		trans := &http.Transport{}
		htclient := &http.Client{Transport: trans}
		response, _ := htclient.Post(url, "text/json", strings.NewReader(""))
		defer response.Body.Close()
		contents, _ := ioutil.ReadAll(response.Body)
		var jsonContents map[string]interface{}
		json.Unmarshal(contents, &jsonContents)
		if response.StatusCode/100 != 2 {
			responseStatus := strconv.Itoa(response.StatusCode)
			contentError := responseStatus + fmt.Sprintf(": %s", jsonContents["error"])
			err = errors.New(contentError)
		} else {
			token["this"] = fmt.Sprintf("%s", contents)
			err = nil
		}
	}
	return token, err
}
