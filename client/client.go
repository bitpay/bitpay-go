//The client package provides convenience methods for authenticating with Bitpay and creating basic invoices.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/vinovest/bitpay-go/pkg/config"
	ku "github.com/vinovest/bitpay-go/pkg/key_utils"
)

type TokenCreation struct {
	Data []struct {
		Policies []struct {
			Policy string   `json:"policy"`
			Method string   `json:"method"`
			Params []string `json:"params"`
		} `json:"policies"`
		Token             string `json:"token"`
		Facade            string `json:"facade"`
		DateCreated       int64  `json:"dateCreated"`
		PairingExpiration int64  `json:"pairingExpiration"`
		PairingCode       string `json:"pairingCode"`
	} `json:"data"`
}

// The Client struct maintains the state of the current client. To use a client from session to session, the Pem and Token will need to be saved and used in the next client. The ClientId can be recreated by using the key_util.GenerateSinFromPem func, and the ApiUri will generally be https://bitpay.com. Insecure should generally be set to false or not set at all, there are a limited number of test scenarios in which it must be set to true.
type Client struct {
	config config.BitpayData
	facade config.Facade
}

func New(c config.BitpayData, f config.Facade) *Client {
	return &Client{
		config: c,
		facade: f,
	}
}

// CreateInvoice returns an invoice type or pass the error from the server. The method will create an invoice on the BitPay server.
func (c *Client) CreateInvoice(payload InvoiceCreation) (Invoice, error) {
	response, _ := c.post("invoices", payload)
	i := Invoice{}
	err := handleResponse(response, &i)
	return i, err
}

func (c *Client) GetTokens() ([]map[string]string, error) {
	url := c.config.GetEnvURL() + "/tokens"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	key, err := c.config.GetPrivateKey(c.facade)
	if err != nil {
		return nil, err
	}
	publ := ku.ExtractCompressedPublicKey(key)

	req.Header.Add("X-Identity", publ)
	sig := ku.Sign(url, key)
	req.Header.Add("X-Signature", sig)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("Did not get ")
	}

	var jsonContents map[string]interface{}
	if err := json.Unmarshal(contents, &jsonContents); err != nil {
		return nil, err
	}

	this, err := json.Marshal(jsonContents["data"])
	if err != nil {
		return nil, err
	}
	var t []map[string]string
	return t, json.Unmarshal(this, &t)
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

func PairClient(key string, env config.Environment, f config.Facade) (TokenCreation, error) {
	sin := ku.GenerateSinFromPem(key)
	url := config.UrlMapping[env] + "/tokens"

	data := map[string]string{
		"id":     sin,
		"facade": string(f),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return TokenCreation{}, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return TokenCreation{}, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenCreation{}, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return TokenCreation{}, nil
	}
	if response.StatusCode/100 != 2 {
		return TokenCreation{}, fmt.Errorf("Did not get status code 200. %v", string(contents))
	}

	t := TokenCreation{}
	err = json.Unmarshal(contents, &t)
	return t, err
}

func (c *Client) post(path string, body interface{}) (*http.Response, error) {
	url := config.UrlMapping[c.config.BitPayConfiguration.Environment] + "/" + path
	payload, err := json.Marshal(body)
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

	key, err := c.config.GetPrivateKey(c.facade)
	if err != nil {
		return nil, err
	}

	publ := ku.ExtractCompressedPublicKey(key)
	req.Header.Add("X-Identity", publ)

	sig := ku.Sign(url+string(payload), key)
	req.Header.Add("X-Signature", sig)

	return http.DefaultClient.Do(req)
}

// GetInvoice is a public facade method, any client which has the ApiUri field set can retrieve an invoice from that endpoint, provided they have the invoice id.
func (c *Client) GetInvoice(invId string) (Invoice, error) {
	url := c.config.GetEnvURL() + "/invoices/" + invId

	response, err := http.Get(url)
	if err != nil {
		return Invoice{}, nil
	}
	i := Invoice{}
	return i, handleResponse(response, &i)
}

func handleResponse(resp *http.Response, i interface{}) error {
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return errors.New("value is not of pointer type")
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("%v", string(contents))
	}

	w := responseWrapper{}
	if err := json.Unmarshal(contents, &w); err != nil {
		return err
	}

	return json.Unmarshal(w.Data, i)
}
