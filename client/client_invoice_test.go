package client_test

import (
	. "github.com/bitpay/bitpay-go/client"
	ku "github.com/bitpay/bitpay-go/key_utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"time"
)

var _ = Describe("CreateInvoice", func() {
	It("creates an invoice for the price and currency sent", func() {
		time.Sleep(5)
		pm := ku.GeneratePem()
		apiuri := os.ExpandEnv("$RCROOTADDRESS")
		webClient := Client{ApiUri: apiuri, Insecure: true, Pem: pm}
		code := os.ExpandEnv("$INVOICEPAIR")
		token, _ := webClient.PairWithCode(code)
		webClient.Token = token
		response, _ := webClient.CreateInvoice(10, "USD")
		Expect(response.Price).To(Equal(10.00))
		response, _ = webClient.CreateInvoice(0.00023, "BTC")
		Expect(response.Price).To(Equal(0.00023))
	})
})
