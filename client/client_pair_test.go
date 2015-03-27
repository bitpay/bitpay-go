package client_test

import (
	. "github.com/bitpay/bitpay-go/client"
	ku "github.com/bitpay/bitpay-go/key_utils"

	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientPair", func() {
	It("pairs with the server with a pairing code", func() {
		time.Sleep(5)
		pm := ku.GeneratePem()
		code := os.ExpandEnv("$PAIRINGCODE")
		apiuri := os.ExpandEnv("$RCROOTADDRESS")
		webClient := Client{ApiUri: apiuri, Insecure: true, Pem: pm}
		token, _ := webClient.PairWithCode(code)
		webClient.Token = token
		Expect(webClient.Token.Facade).To(Equal("pos"))
	})
})
