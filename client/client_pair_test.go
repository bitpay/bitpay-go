package client_test

import (
	. "github.com/bitpay/bitpay-go/client"
	ku "github.com/bitpay/bitpay-go/key_utils"

	"io/ioutil"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientPair", func() {
	It("pairs with the server with a pairing code", func() {
		time.Sleep(30)
		pm := ku.GeneratePem()
		apiuri := os.ExpandEnv("$RCROOTADDRESS")
		webClient := Client{ApiUri: apiuri, Insecure: true, Pem: pm}
		gopath := os.ExpandEnv("$GOPATH")
		tempFolder := gopath + "/temp/"
		code, err := ioutil.ReadFile(tempFolder + "paircode.txt")
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
		Expect(webClient.Token.Facade).To(Equal("pos"))
	})
})
