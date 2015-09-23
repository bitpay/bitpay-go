package client_test

import (
	"encoding/json"
	. "github.com/bitpay/bitpay-go/client"
	ku "github.com/bitpay/bitpay-go/key_utils"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("ClientPair", func() {
	It("pairs with the server with a pairing code", func() {
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
			Expect(res.StatusCode).To(Equal(200), string(contents)+"token: "+mertok)
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
		Expect(webClient.Token.Facade).To(Equal("pos"))
	})
})

func processToken(response *http.Response, jsonContents map[string]interface{}) (tok Token, err error) {
	datarray := jsonContents["data"].([]interface{})
	data, _ := json.Marshal(datarray[0])
	json.Unmarshal(data, &tok)
	return tok, nil
}
