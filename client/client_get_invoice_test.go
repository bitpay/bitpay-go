package client_test

import (
	. "github.com/bitpay/bitpay-go/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"io/ioutil"
	"os"
	"time"
)

var _ = Describe("RetrieveInvoice", func() {
	It("Retrieves an invoice from the server with an id", func() {
		time.Sleep(3 * time.Second)
		pm := os.ExpandEnv("$BITPAYPEM")
		apiuri := os.ExpandEnv("$BITPAYAPI")
		webClient := Client{ApiUri: apiuri, Insecure: true, Pem: pm}
		invoiceId := os.ExpandEnv("$INVOICEID")
		retrievedInvoice, err := webClient.GetInvoice(invoiceId)
		if err != nil {
			println("\n" + webClient.ApiUri + " errored retrieving an invoice: Error - " + err.Error())
		}
		Expect(retrievedInvoice.Id).To(Equal(invoiceId))
	})
})
