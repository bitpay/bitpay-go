//The client package provides convenience methods for authenticating with Bitpay and creating basic invoices.
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	ku "github.com/bitpay/bitpay-go/key_utils"
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

type Invoice struct {
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
func (client *Client) CreateInvoice(price float64, currency string) (*Invoice, error) {
	match, _ := regexp.MatchString("^[[:upper:]]{3}$", currency)
	if !match {
		return nil, errors.New("BitPayArgumentError: invalid currency code")
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
	response, err := client.Post("invoices", paylo)
	if err != nil {
		return nil, err
	}
	return processInvoice(response)
}

// PairWithFacade
func (client *Client) PairWithFacade(str string) (*Token, error) {
	paylo := make(map[string]string)
	paylo["facade"] = str
	return client.PairClient(paylo)
}

// PairWithCode retrieves a token from the server and authenticates the keys of the calling client. The string passed to the client is a "pairing code" that must be retrieved from https://bitpay.com/dashboard/merchant/api-tokens. PairWithCode returns a Token type that must be assigned to the Token field of a client in order for that client to create invoices. For example `client.Token = client.PairWithCode("abcdefg")`.
func (client *Client) PairWithCode(str string) (*Token, error) {
	match, _ := regexp.MatchString("^[[:alnum:]]{7}$", str)
	if !match {
		return nil, errors.New("BitPayArgumentError: invalid pairing code")
	}
	paylo := make(map[string]string)
	paylo["pairingCode"] = str
	return client.PairClient(paylo)
}

func (client *Client) PairClient(paylo map[string]string) (*Token, error) {
	paylo["id"] = client.ClientId
	sin := ku.GenerateSinFromPem(client.Pem)
	client.ClientId = sin
	url := client.ApiUri + "/tokens"
	htclient := setHttpClient(client)
	payload, err := json.Marshal(paylo)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	response, err := htclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var jsonContents map[string]interface{}
	if err = json.Unmarshal(contents, &jsonContents); err != nil {
		return nil, err
	}
	if response.StatusCode/100 != 2 {
		return nil, processErrorMessage(response, jsonContents)
	}
	return processToken(response, jsonContents)
}

func (client *Client) Post(path string, paylo map[string]string) (*http.Response, error) {
	url := client.ApiUri + "/" + path
	htclient := setHttpClient(client)
	payload, err := json.Marshal(paylo)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	publ := ku.ExtractCompressedPublicKey(client.Pem)
	req.Header.Add("X-Identity", publ)
	sig := ku.Sign(url+string(payload), client.Pem)
	req.Header.Add("X-Signature", sig)
	return htclient.Do(req)
}

// GetInvoice is a public facade method, any client which has the ApiUri field set can retrieve an invoice from that endpoint, provided they have the invoice id.
func (client *Client) GetInvoice(invId string) (*Invoice, error) {
	url := client.ApiUri + "/invoices/" + invId
	htclient := setHttpClient(client)
	response, err := htclient.Get(url)
	if err != nil {
		return nil, err
	}
	return processInvoice(response)
}

func (client *Client) GetTokens() (tokes []map[string]string, err error) {
	url := client.ApiUri + "/tokens"
	htclient := setHttpClient(client)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	publ := ku.ExtractCompressedPublicKey(client.Pem)
	req.Header.Add("X-Identity", publ)
	sig := ku.Sign(url, client.Pem)
	req.Header.Add("X-Signature", sig)
	response, err := htclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var jsonContents map[string]interface{}
	if err = json.Unmarshal(contents, &jsonContents); err != nil {
		return nil, err
	}
	if response.StatusCode/100 != 2 {
		return nil, processErrorMessage(response, jsonContents)
	}
	this, err := json.Marshal(jsonContents["data"])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(this, &tokes)
	return tokes, err
}

func (client *Client) GetToken(facade string) (token string, err error) {
	tokens, err := client.GetTokens()
	if err != nil {
		return "", err
	}
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

func processToken(response *http.Response, jsonContents map[string]interface{}) (*Token, error) {
	datarray := jsonContents["data"].([]interface{})
	data, err := json.Marshal(datarray[0])
	if err != nil {
		return nil, err
	}
	tok := new(Token)
	err = json.Unmarshal(data, tok)
	return tok, err
}

func processInvoice(response *http.Response) (*Invoice, error) {
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var jsonContents map[string]interface{}
	if err := json.Unmarshal(contents, &jsonContents); err != nil {
		return nil, err
	}
	if response.StatusCode/100 != 2 {
		return nil, processErrorMessage(response, jsonContents)
	}
	this, err := json.Marshal(jsonContents["data"])
	if err != nil {
		return nil, err
	}
	inv := new(Invoice)
	err = json.Unmarshal(this, inv)
	return inv, err
}
