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
	"time"
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

type Invoice struct {
	Addresses            *PaymentDisplay `json:"addresses,omitempty"`
	AmountPaid           int             `json:"amountPaid,omitempty"`
	BitcoinAddress       string          `json:"bitcoinAddress,omitempty"`
	BtcDue               string          `json:"btcDue,omitempty"`
	BtcPaid              string          `json:"btcPaid,omitempty"`
	BtcPrice             string          `json:"btcPrice,omitempty"`
	Buyer                *Buyer          `json:"buyer,omitempty"`
	BuyerPaidBtcMinerFee interface{}     `json:"buyerPaidBtcMinerFee,omitempty"`
	BuyerTotalBtcAmount  interface{}     `json:"buyerTotalBtcAmount,omitempty"`
	CryptoInfo           []struct {
		Address    string `json:"address,omitempty"`
		CryptoCode string `json:"cryptoCode,omitempty"`
		CryptoPaid string `json:"cryptoPaid,omitempty"`
		Due        string `json:"due,omitempty"`
		ExRates    struct {
			USD int `json:"USD,omitempty"`
		} `json:"exRates,omitempty"`
		NetworkFee  string `json:"networkFee,omitempty"`
		Paid        string `json:"paid,omitempty"`
		PaymentType string `json:"paymentType,omitempty"`
		PaymentUrls struct {
			BIP21  string      `json:"BIP21,omitempty"`
			BIP72  interface{} `json:"BIP72,omitempty"`
			BIP72B interface{} `json:"BIP72b,omitempty"`
			BIP73  interface{} `json:"BIP73,omitempty"`
			BOLT11 interface{} `json:"BOLT11,omitempty"`
		} `json:"paymentUrls,omitempty"`
		Payments []interface{} `json:"payments,omitempty"`
		Price    string        `json:"price,omitempty"`
		Rate     float64       `json:"rate,omitempty"`
		TotalDue string        `json:"totalDue,omitempty"`
		TxCount  int           `json:"txCount,omitempty"`
		URL      string        `json:"url,omitempty"`
	} `json:"cryptoInfo,omitempty"`
	Currency    string `json:"currency,omitempty"`
	CurrentTime int64  `json:"currentTime,omitempty"`
	/*
		ExRates     struct {
			USD int `json:"USD,omitempty"`
		} `json:"exRates,omitempty"`
	*/

	ExceptionStatus bool `json:"exceptionStatus,omitempty"`
	/*
		ExchangeRates   struct {
			BTC struct {
				USD int `json:"USD,omitempty"`
			} `json:"BTC,omitempty"`
		} `json:"exchangeRates,omitempty"`
	*/

	ExpirationTime int64 `json:"expirationTime,omitempty"`
	/*
		Flags          struct {
			Refundable bool `json:"refundable,omitempty"`
		} `json:"flags,omitempty"`
	*/
	GUID           string      `json:"guid,omitempty"`
	ID             string      `json:"id,omitempty"`
	InvoiceTime    int64       `json:"invoiceTime,omitempty"`
	ItemCode       interface{} `json:"itemCode,omitempty"`
	ItemDesc       interface{} `json:"itemDesc,omitempty"`
	LowFeeDetected bool        `json:"lowFeeDetected,omitempty"`
	/*MinerFees      struct {
		BTC struct {
			SatoshisPerByte int `json:"satoshisPerByte,omitempty"`
			TotalFee        int `json:"totalFee,omitempty"`
		} `json:"BTC,omitempty"`
	} `json:"minerFees,omitempty"`
	*/
	OrderID interface{} `json:"orderId,omitempty"`
	/*PaymentCodes struct {
		BTC struct {
			BIP21  string      `json:"BIP21,omitempty"`
			BIP72  interface{} `json:"BIP72,omitempty"`
			BIP72B interface{} `json:"BIP72b,omitempty"`
			BIP73  interface{} `json:"BIP73,omitempty"`
			BOLT11 interface{} `json:"BOLT11,omitempty"`
		} `json:"BTC,omitempty"`
	} `json:"paymentCodes,omitempty"`
	*/
	/*PaymentSubtotals struct {
		BTC int `json:"BTC,omitempty"`
	} `json:"paymentSubtotals,omitempty"`
	PaymentTotals struct {
		BTC int `json:"BTC,omitempty"`
	} `json:"paymentTotals,omitempty"`
	PaymentUrls struct {
		BIP21  string      `json:"BIP21,omitempty"`
		BIP72  interface{} `json:"BIP72,omitempty"`
		BIP72B interface{} `json:"BIP72b,omitempty"`
		BIP73  interface{} `json:"BIP73,omitempty"`
		BOLT11 interface{} `json:"BOLT11,omitempty"`
	} `json:"paymentUrls,omitempty"`
	*/
	PosData                     interface{} `json:"posData,omitempty"`
	Price                       float64     `json:"price,omitempty"`
	Rate                        float64     `json:"rate,omitempty"`
	RefundAddressRequestPending bool        `json:"refundAddressRequestPending,omitempty"`
	Status                      string      `json:"status,omitempty"`
	/*SupportedTransactionCurrencies struct {
		BTC struct {
			Enabled bool        `json:"enabled"`
			Reason  interface{} `json:"reason"`
		} `json:"BTC"`
	} `json:"supportedTransactionCurrencies,omitempty"`
	*/
	PaymentCurrencies []string `json:"paymentCurrencies,omitempty"`
	Token             string   `json:"token,omitempty"`
	URL               string   `json:"url,omitempty"`
}

type PaymentDisplay struct {
	Btc  string `json:"BTC,omitempty"`
	Bch  string `json:"BCH,omitempty"`
	Eth  string `json:"ETH,omitempty"`
	Gusd string `json:"GUSD,omitempty"`
	Pax  string `json:"PAX,omitempty"`
	Busd string `json:"BUSD,omitempty"`
	Usdc string `json:"USDC,omitempty"`
	Xrp  string `json:"XRP,omitempty"`
}

type Buyer struct {
	Address1   interface{} `json:"address1,omitempty"`
	Address2   interface{} `json:"address2,omitempty"`
	Country    interface{} `json:"country,omitempty"`
	Email      interface{} `json:"email,omitempty"`
	Locality   interface{} `json:"locality,omitempty"`
	Name       interface{} `json:"name,omitempty"`
	Phone      interface{} `json:"phone,omitempty"`
	PostalCode interface{} `json:"postalCode,omitempty"`
	Region     interface{} `json:"region,omitempty"`
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

type Payout struct {
	Id                string
	Account           string
	Reference         string
	SupportPhone      string
	Status            string
	Amount            float64
	PercentFee        float64
	Fee               float64
	DepositTotal      float64
	Btc               float64
	Currency          string
	RequestDate       time.Time
	EffectiveDate     time.Time
	NotificationUrl   string
	NotificationEmail string
	Instructions      []Instruction
	Token             string
}

type Btc struct {
	Unpaid int
	Paid   int
}

type Instruction struct {
	Id             string
	Amount         float64
	Btc            Btc
	Address        string
	Label          string
	Status         string
	WalletProvider string
	Receiverinfo   Receiverinfo
}

type Receiverinfo struct {
	Name         string
	EmailAddress string
	Address      Address
}

type Address struct {
	StreetAddress1 string
	StreetAddress2 string
	Locality       string
	Region         string
	PostalCode     string
	Country        string
}

// CreateInvoice returns an invoice type or pass the error from the server. The method will create an invoice on the BitPay server.
func (client *Client) CreateInvoice(i Invoice) (inv Invoice, err error) {
	match, _ := regexp.MatchString("^[[:upper:]]{3}$", i.Currency)
	if !match {
		err = errors.New("BitPayArgumentError: invalid currency code")
		return inv, err
	}

	i.Token = client.Token.Token
	response, _ := client.Post("invoices", i)
	var invoice Invoice
	invoice, err = processInvoice(response)
	return invoice, err
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

func (client *Client) Post(path string, paylo interface{}) (response *http.Response, err error) {
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
func (client *Client) GetInvoice(invId string) (inv Invoice, err error) {
	url := client.ApiUri + "/invoices/" + invId
	htclient := setHttpClient(client)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-accept-version", "2.0.0")
	publ := ku.ExtractCompressedPublicKey(client.Pem)
	req.Header.Add("X-Identity", publ)
	sig := ku.Sign(url, client.Pem)
	req.Header.Add("X-Signature", sig)
	response, _ := htclient.Do(req)
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
	response, err := htclient.Do(req)
	if err != nil {
		return tokes, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return tokes, err
	}
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
	var contentError string
	if jsonContents != nil {
		contentError = responseStatus + ": " + jsonContents["error"].(string)
	} else {
		contentError = responseStatus
	}
	return errors.New(contentError)
}

func processToken(response *http.Response, jsonContents map[string]interface{}) (tok Token, err error) {
	datarray := jsonContents["data"].([]interface{})
	data, _ := json.Marshal(datarray[0])
	json.Unmarshal(data, &tok)
	return tok, nil
}

func processInvoice(response *http.Response) (inv Invoice, err error) {
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

// CreatePayout create and returns a PayOut.
func (client *Client) CreatePayout(p Payout) (payout Payout, err error) {
	match, _ := regexp.MatchString("^[[:upper:]]{3}$", p.Currency)
	if !match {
		err = errors.New("BitPayArgumentError: invalid currency code")
		return payout, err
	}
	p.Token = client.Token.Token
	response, err := client.Post("payouts", p)
	body, err := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &payout)
	return payout, err
}
