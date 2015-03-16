package client_test

import (
	. "bitpay/client"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os/exec"
)

var _ = Describe("ClientPair", func() {
	It("pairs with the server with a pairing code", func() {
		gopath := os.ExpandEnv("$GOPATH")
		pyloc := gopath + "/helpers/pair_steps.py"
		cmd := exec.Command(pyloc)
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		byt, _ := ioutil.ReadAll(stdout)
		code := string(byt)
		apiuri := os.ExpandEnv("$RCROOTADDRESS")
		webClient := Client{ApiUri: apiuri, Insecure: true}
		token, _ := webClient.PairWithCode(code)
		Expect(token["facade"]).To(Equal("pos"))
	})
})
