package client

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type InvoiceCreation struct {
	// These are required fields
	Currency string  `json:"currency"`
	Price    float64 `json:"price"`
	Token    string  `json:"token"`

	// These are optional fields below
	OrderID               string `json:"orderId"`
	FullNotifications     bool   `json:"fullNotifications,omitempty"`
	ExtendedNotifications bool   `json:"extendedNotifications,omitempty"`
	TransactionSpeed      string `json:"transactionSpeed,omitempty"`
	NotificationURL       string `json:"notificationURL,omitempty"`
	NotificationEmail     string `json:"notificationEmail,omitempty"`
	RedirectURL           string `json:"redirectURL,omitempty"`
	Buyer                 Buyer  `json:"buyer,omitempty"`
	PosData               string `json:"posData,omitempty"`
	ItemDesc              string `json:"itemDesc,omitempty"`
}

func (i InvoiceCreation) validate() error {
	if i.Currency == "" {
		return errors.New("Need to specify a currency")
	}

	if i.Price <= 0 {
		return errors.New("Invalid price")
	}

	if i.Token == "" {
		return errors.New("No token set")
	}

	return nil
}

type InvoicePayload struct {
	Price    string `json:"price"`
	Currency string `json:"currency"`
	Token    string `json:"token"`
}

type responseWrapper struct {
	Facade string          `json:"facade"`
	Data   json.RawMessage `json:"data"`
}

type ExRates struct {
	ETH  decimal.Decimal `json:"ETH"`
	EUR  decimal.Decimal `json:"EUR"`
	BTC  decimal.Decimal `json:"BTC"`
	USD  decimal.Decimal `json:"USD"`
	BCH  decimal.Decimal `json:"BCH"`
	GUSD decimal.Decimal `json:"GUSD"`
	PAX  decimal.Decimal `json:"PAX"`
	BUSD decimal.Decimal `json:"BUSD"`
	USDC decimal.Decimal `json:"USDC"`
	XRP  decimal.Decimal `json:"XRP"`
}
type Transactions struct {
	Amount        int64     `json:"amount"`
	Confirmations int       `json:"confirmations"`
	Time          time.Time `json:"time"`
	ReceivedTime  time.Time `json:"receivedTime"`
	Txid          string    `json:"txid"`
	ExRates       ExRates   `json:"exRates"`
	Type          string    `json:"type,omitempty"`
	RefundAmount  int64     `json:"refundAmount,omitempty"`
}

type EnabledCurrency struct {
	Enabled bool `json:"enbled"`
}

type Buyer struct {
	Name       string `json:"name"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	Locality   string `json:"locality"`
	Region     string `json:"region"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Notify     bool   `json:"notify"`
}

type RefundAddressData struct {
	Type string      `json:"type"`
	Date time.Time   `json:"date"`
	Tag  interface{} `json:"tag"`
}

type BuyerProvidedInfo struct {
	Name                        string `json:"name"`
	PhoneNumber                 string `json:"phoneNumber"`
	SelectedWallet              string `json:"selectedWallet"`
	SelectedTransactionCurrency string `json:"selectedTransactionCurrency"`
	EmailAddress                string `json:"emailAddress"`
}

type MinerFee struct {
	SatoshisPerByte int `json:"satoshisPerByte"`
	TotalFee        int `json:"totalFee"`
}

type Shopper struct {
	User string `json:"user"`
}

type RefundInfo struct {
	SupportRequest string             `json:"supportRequest"`
	Currency       string             `json:"currency"`
	Amounts        map[string]float64 `json:"amounts"` // in fiat and crypto
}

type CryptoCurrency string

const (
	BTC  CryptoCurrency = "BTC"
	BCH  CryptoCurrency = "BCH"
	ETH  CryptoCurrency = "ETH"
	GUSD CryptoCurrency = "GUSD"
	PAX  CryptoCurrency = "PAX"
	BUSD CryptoCurrency = "BUSD"
	USDC CryptoCurrency = "USDC"
	XRP  CryptoCurrency = "XRP"
)

type Invoice struct {
	URL                            string                                        `json:"url"`
	PosData                        string                                        `json:"posData"`
	Status                         string                                        `json:"status"`
	Price                          int                                           `json:"price"`
	Currency                       string                                        `json:"currency"`
	ItemDesc                       string                                        `json:"itemDesc"`
	OrderID                        string                                        `json:"orderId"`
	InvoiceTime                    int64                                         `json:"invoiceTime"`
	ExpirationTime                 int64                                         `json:"expirationTime"`
	CurrentTime                    int64                                         `json:"currentTime"`
	ID                             string                                        `json:"id"`
	LowFeeDetected                 bool                                          `json:"lowFeeDetected"`
	AmountPaid                     int64                                         `json:"amountPaid"`
	DisplayAmountPaid              string                                        `json:"displayAmountPaid"`
	ExceptionStatus                bool                                          `json:"exceptionStatus"`
	TargetConfirmations            int                                           `json:"targetConfirmations"`
	Transactions                   []Transactions                                `json:"transactions"`
	TransactionSpeed               string                                        `json:"transactionSpeed"`
	Buyer                          Buyer                                         `json:"buyer"`
	RedirectURL                    string                                        `json:"redirectURL"`
	RefundAddresses                []map[string]RefundAddressData                `json:"refundAddresses"`
	RefundAddressRequestPending    bool                                          `json:"refundAddressRequestPending"`
	BuyerProvidedEmail             string                                        `json:"buyerProvidedEmail"`
	BuyerProvidedInfo              BuyerProvidedInfo                             `json:"buyerProvidedInfo"`
	PaymentSubtotals               map[CryptoCurrency]decimal.Decimal            `json:"paymentSubtotals"`
	PaymentTotals                  map[CryptoCurrency]decimal.Decimal            `json:"paymentTotals"`
	PaymentDisplayTotals           map[CryptoCurrency]decimal.Decimal            `json:"paymentDisplayTotals"`
	PaymentDisplaySubTotals        map[CryptoCurrency]decimal.Decimal            `json:"paymentDisplaySubTotals"`
	ExchangeRates                  map[CryptoCurrency]map[string]decimal.Decimal `json:"exchangeRates"`
	MinerFees                      map[CryptoCurrency]MinerFee                   `json:"minerFees"`
	Shopper                        Shopper                                       `json:"shopper"`
	JSONPayProRequired             bool                                          `json:"jsonPayProRequired"`
	RefundInfo                     []RefundInfo                                  `json:"refundInfo"`
	TransactionCurrency            string                                        `json:"transactionCurrency"`
	SupportedTransactionCurrencies map[CryptoCurrency]EnabledCurrency            `json:"supportedTransactionCurrencies"`
	PaymentCodes                   map[CryptoCurrency]map[string]string          `json:"paymentCodes"`
	Token                          string                                        `json:"token"`
}
