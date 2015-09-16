//The client package provides convenience methods for authenticating with Bitpay and creating basic invoices.
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	ku "github.com/bitpay/bitpay-go/key_utils"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

// The Client struct maintains the state of the current client. To use a client from session to session, the Pem and Token will need to be saved and used in the next client. The ClientId can be recreated by using the key_util.GenerateSinFromPem func, and the ApiUri will generally be https://bitpay.com. Insecure should generally be set to false or not set at all, there are a limited number of test scenarios in which it must be set to true.
type Client struct {
	Pem      string
	ApiUri   string
	Insecure bool
	ClientId string
	Token    Token
}

// The Token struct is a go mapping of a subset of the JSON returned from the server with a request to authenticate (pair).
type Token struct {
	Token             string
	Facade            string
	DateCreated       float64
	PairingExpiration float64
	Resource          string
	PairingCode       string
}

// Go struct mapping the JSON returned from the BitPay server when sending a POST or GET request to /invoices.

type invoice struct {
	Url             string
	Status          string
	BtcPrice        string
	BtcDue          string
	Price           float64
	Currency        string
	ExRates         map[string]float64
	InvoiceTime     int64
	ExpirationTime  int64
	CurrentTime     int64
	Guid            string
	Id              string
	BtcPaid         string
	Rate            float64
	ExceptionStatus bool
	PaymentUrls     map[string]string
	Token           string
}

// CreateInvoice returns an invoice type or pass the error from the server. The method will create an invoice on the BitPay server.
func (client *Client) CreateInvoice(price float64, currency string) (inv invoice, err error) {
	match, _ := regexp.MatchString("^[[:upper:]]{3}$", currency)
	if !match {
		err = errors.New("BitPayArgumentError: invalid currency code")
		return inv, err
	}
	paylo := make(map[string]string)
	var floatPrec int
	if currency == "BTC" {
		floatPrec = 8
	} else {
		floatPrec = 2
	}
	priceString := strconv.FormatFloat(price, 'f', floatPrec, 64)
	paylo["price"] = priceString
	paylo["currency"] = currency
	paylo["token"] = client.Token.Token
	paylo["id"] = client.ClientId
	response, _ := client.Post("invoices", paylo)
	inv, err = processInvoice(response)
	return inv, err
}

// PairWithFacade
func (client *Client) PairWithFacade(str string) (tok Token, err error) {
	paylo := make(map[string]string)
	paylo["facade"] = str
	tok, err = client.PairClient(paylo)
	return tok, err
}

// PairWithCode retrieves a token from the server and authenticates the keys of the calling client. The string passed to the client is a "pairing code" that must be retrieved from https://bitpay.com/dashboard/merchant/api-tokens. PairWithCode returns a Token type that must be assigned to the Token field of a client in order for that client to create invoices. For example `client.Token = client.PairWithCode("abcdefg")`.
func (client *Client) PairWithCode(str string) (tok Token, err error) {
	match, _ := regexp.MatchString("^[[:alnum:]]{7}$", str)
	if !match {
		err = errors.New("BitPayArgumentError: invalid pairing code")
		return tok, err
	}
	paylo := make(map[string]string)
	paylo["pairingCode"] = str
	tok, err = client.PairClient(paylo)
	return tok, err
}

func (client *Client) PairClient(paylo map[string]string) (tok Token, err error) {
	paylo["id"] = client.ClientId
	sin := ku.GenerateSinFromPem(client.Pem)
	client.ClientId = sin
	url := client.ApiUri + "/tokens"
	htclient := setHttpClient(client)
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
		tok, err = processToken(response, jsonContents)
		err = nil
	}
	return tok, err
}

func (client *Client) Post(path string, paylo map[string]string) (response *http.Response, err error) {
	url := client.ApiUri + "/" + path
	htclient := setHttpClient(client)
	payload, _ := json.Marshal(paylo)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	publ := ku.ExtractCompressedPublicKey(client.Pem)
	req.Header.Add("X-Identity", publ)
	sig := ku.Sign(url+string(payload), client.Pem)
	req.Header.Add("X-Signature", sig)
	response, err = htclient.Do(req)
	return response, err
}

// GetInvoice is a public facade method, any client which has the ApiUri field set can retrieve an invoice from that endpoint, provided they have the invoice id.
func (client *Client) GetInvoice(invId string) (inv invoice, err error) {
	url := client.ApiUri + "/invoices/" + invId
	htclient := setHttpClient(client)
	response, _ := htclient.Get(url)
	inv, err = processInvoice(response)
	return inv, err
}

func (client *Client) GetTokens() (tokes []map[string]string, err error) {
	url := client.ApiUri + "/tokens"
	htclient := setHttpClient(client)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	publ := ku.ExtractCompressedPublicKey(client.Pem)
	req.Header.Add("X-Identity", publ)
	sig := ku.Sign(url, client.Pem)
	req.Header.Add("X-Signature", sig)
	response, _ := htclient.Do(req)
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	var jsonContents map[string]interface{}
	json.Unmarshal(contents, &jsonContents)
	if response.StatusCode/100 != 2 {
		err = processErrorMessage(response, jsonContents)
	} else {
		this, _ := json.Marshal(jsonContents["data"])
		json.Unmarshal(this, &tokes)
		err = nil
	}
	return tokes, nil
}

func (client *Client) GetToken(facade string) (token string, err error) {
	tokens, _ := client.GetTokens()
	for _, token := range tokens {
		toke, ok := token[facade]
		if ok {
			return toke, nil
		}
	}
	return "error", errors.New("facade not available in tokens")
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

func processToken(response *http.Response, jsonContents map[string]interface{}) (tok Token, err error) {
	datarray := jsonContents["data"].([]interface{})
	data, _ := json.Marshal(datarray[0])
	json.Unmarshal(data, &tok)
	return tok, nil
}

func processInvoice(response *http.Response) (inv invoice, err error) {
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	var jsonContents map[string]interface{}
	json.Unmarshal(contents, &jsonContents)
	if response.StatusCode/100 != 2 {
		err = processErrorMessage(response, jsonContents)
	} else {
		this, _ := json.Marshal(jsonContents["data"])
		json.Unmarshal(this, &inv)
		err = nil
	}
	return inv, err
}
