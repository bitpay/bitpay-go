package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

var UrlMapping = map[Environment]string{
	Test:       "https://test.bitpay.com",
	Production: "https://bitpay.com",
}

type Environment string

const (
	Production Environment = "Prod"
	Test       Environment = "Test"
)

type Facade string

const (
	Merchant Facade = "merchant"
	Payroll  Facade = "payroll"
	Public   Facade = "public"
)

type EnvironmentData struct {
	PrivateKeyPath string            `json:"PrivateKeyPath"`
	PrivateKey     string            `json:"PrivateKey"`
	APITokens      map[Facade]string `json:"ApiTokens"`
}

type BitpayData struct {
	BitPayConfiguration struct {
		Environment Environment                     `json:"Environment"`
		EnvConfig   map[Environment]EnvironmentData `json:"EnvConfig"`
	} `json:"BitPayConfiguration"`
}

func (b BitpayData) GetEnvURL() string {
	return UrlMapping[b.BitPayConfiguration.Environment]
}

func (b BitpayData) GetPrivateKey(f Facade) (string, error) {
	key := b.BitPayConfiguration.EnvConfig[Test].PrivateKey
	path := b.BitPayConfiguration.EnvConfig[Test].PrivateKeyPath
	if v, ok := b.BitPayConfiguration.EnvConfig[Production]; ok {
		key = v.PrivateKey
		path = v.PrivateKeyPath
	}

	if key != "" {
		return key, nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "Unable to load file")
	}

	return string(data), nil
}

func (b BitpayData) GetToken(f Facade) string {
	if v, ok := b.BitPayConfiguration.EnvConfig[Production]; ok {
		return v.APITokens[f]
	}
	return b.BitPayConfiguration.EnvConfig[Test].APITokens[f]
}

func LoadFromString(s string) (BitpayData, error) {
	b := BitpayData{}
	err := json.Unmarshal([]byte(s), &b)
	return b, err
}

func LoadFromFile(path string) (BitpayData, error) {

	f, err := os.Open(path)
	if err != nil {
		return BitpayData{}, err
	}

	b := BitpayData{}
	err = json.NewDecoder(f).Decode(&b)
	return b, err
}
