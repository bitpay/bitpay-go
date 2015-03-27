package client_test

import (
	. "github.com/bitpay/bitpay-go/client"
	ku "github.com/bitpay/bitpay-go/key_utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"time"
)

var _ = Describe("CreateInvoice", func() {
	It("creates an invoice for the price and currency sent", func() {
		time.Sleep(30)
		pm := ku.GeneratePem()
		apiuri := os.ExpandEnv("$RCROOTADDRESS")
		webClient := Client{ApiUri: apiuri, Insecure: true, Pem: pm}
		gopath := os.ExpandEnv("$GOPATH")
		tempFolder := gopath + "/temp/"
		code, err := ioutil.ReadFile(tempFolder + "invoicecode.txt")
		if err != nil {
			println(err.Error())
		} else {
			println(code)
		}
		token, err := webClient.PairWithCode(string(code))
		if err != nil {
			println(err.Error())
		} else {
			println(token.Token)
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
