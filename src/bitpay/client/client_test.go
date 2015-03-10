package client_test

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"

	. "bitpay/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Describe("PairWithCode", func() {
		It("Short Circuits on Invalid Pairing Code", func() {
			client := new(Client)
			_, err := client.PairWithCode("abcdefgh")
			Ω(err).Should(MatchError("BitPayArgumentError: invalid pairing code"))
		})
		It("Handles Errors Gracefully", func() {
			server := httptest.NewServer(testHandlers())
			client := Client{ApiUri: server.URL}
			_, err := client.PairWithCode("abcdefg")
			Ω(err).Should(MatchError("407: This error is fake"))
		})
	})
})

//handlers function
//each handler function

func testHandlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/tokens", createTokenHandler).Methods("POST")
	return r
}

func createTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(407)
	w.Write([]byte("{\n  \"error\": \"This error is fake\"\n}"))
}
