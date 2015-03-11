package client_test

import (
	. "bitpay/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os/exec"
)

var _ = Describe("ClientPair", func() {
	//	var page *agouti.Page
	//
	//	BeforeEach(func() {
	//		var err error
	//		page, err = agoutiDriver.NewPage()
	//		Expect(err).NotTo(HaveOccurred())
	//	})
	//
	//	AfterEach(func() {
	//		Expect(page.Destroy()).To(Succeed())
	//	})
	//
	It("pairs with the server", func() {
		cmd := exec.Command("echo", "-n", "Qc7AyCw")
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		byt, _ := ioutil.ReadAll(stdout)
		code := string(byt)
		webClient := Client{ApiUri: "https://paul.bp:8088", Insecure: true}
		token, _ := webClient.PairWithCode(code)
		Expect(token["facade"]).To(Equal("pos"))
		Expect(webClient.ApiUri).To(Equal("https://paul.bp:8088"))
	})
})
