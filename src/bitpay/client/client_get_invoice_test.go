package client_test

import (
	. "bitpay/client"
	ku "bitpay/key_utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"os/exec"
)

var _ = Describe("RetrieveInvoice", func() {
	It("Retrieves an invoice from the server with an id", func() {
		pm := ku.GeneratePem()
		gopath := os.ExpandEnv("$GOPATH")
		apiuri := os.ExpandEnv("$RCROOTADDRESS")
		webClient := Client{ApiUri: apiuri, Insecure: true, Pem: pm}
		var code string
		pyloc := gopath + "/helpers/pair_steps.py"
		cmd := exec.Command(pyloc)
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		byt, _ := ioutil.ReadAll(stdout)
		code = string(byt)
		token, _ := webClient.PairWithCode(code)
		webClient.Token = token
		response, _ := webClient.CreateInvoice(10, "USD")
		invoiceId := response.Id
		retrievedInvoice, _ := webClient.GetInvoice(invoiceId)
		Expect(retrievedInvoice.Id).To(Equal(invoiceId))
	})
})
