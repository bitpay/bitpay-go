package client_test

import (
	. "github.com/bitpay/bitpay-go/client"
	ku "github.com/bitpay/bitpay-go/key_utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var _ = Describe("RetrieveInvoice", func() {
	It("Retrieves an invoice from the server with an id", func() {
		time.Sleep(5)
		pm := ku.GeneratePem()
		apiuri := os.ExpandEnv("$RCROOTADDRESS")
		webClient := Client{ApiUri: apiuri, Insecure: true, Pem: pm}
		code := os.ExpandEnv("RETRIEVEPAIR")
		token, _ := webClient.PairWithCode(code)
		webClient.Token = token
		response, _ := webClient.CreateInvoice(10, "USD")
		invoiceId := response.Id
		retrievedInvoice, _ := webClient.GetInvoice(invoiceId)
		Expect(retrievedInvoice.Id).To(Equal(invoiceId))
	})
})
