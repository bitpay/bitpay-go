package client_test

import (
	"encoding/json"
	. "github.com/bitpay/bitpay-go/client"
	ku "github.com/bitpay/bitpay-go/key_utils"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"time"
)

var _ = Describe("CreateInvoice", func() {
	It("creates an invoice for the price and currency sent", func() {
		time.Sleep(3 * time.Second)
		pm := os.ExpandEnv("$BITPAYPEM")
		pm = strings.Replace(pm, "\\n", "\n", -1)
		apiuri := os.ExpandEnv("$BITPAYAPI")
		webClient := Client{ApiUri: apiuri, Insecure: true, Pem: pm, ClientId: ku.GenerateSinFromPem(pm)}
		mertok, _ := webClient.GetToken("merchant")
		params := make(map[string]string)
		params["token"] = mertok
		params["facade"] = "pos"
		params["id"] = webClient.ClientId
		res, err := webClient.Post("tokens", params)
		defer res.Body.Close()
		contents, _ := ioutil.ReadAll(res.Body)
		var jsonContents map[string]interface{}
		json.Unmarshal(contents, &jsonContents)
		var tok Token
		if res.StatusCode/100 != 2 {
			Expect(res.StatusCode).To(Equal(200))
		} else {
			tok, err = processToken(res, jsonContents)
			err = nil
		}
		code := tok.PairingCode
		token, err := webClient.PairWithCode(string(code))
		if err != nil {
			Expect(err.Error).To(Equal("Should be no error"))
		}
		webClient.Token = token
		response, err := webClient.CreateInvoice(10, "USD")
		if err != nil {
			println("The test errored when creating an invoice")
		}
		Expect(response.Price).To(Equal(10.00))
		response, _ = webClient.CreateInvoice(0.00023, "BTC")
		Expect(response.Price).To(Equal(0.00023))
	})
})
