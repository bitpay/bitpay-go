package client_test

import (
	. "github.com/bitpay/bitpay-go/client"
	ku "github.com/bitpay/bitpay-go/key_utils"

	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Describe("PairWithCode", func() {
		It("Short Circuits on Invalid Pairing Code", func() {
			client := new(Client)
			_, err := client.PairWithCode("abcdefgh")
			Expect(err).To(MatchError("BitPayArgumentError: invalid pairing code"))
		})
		It("Handles Errors Gracefully", func() {
			pm := ku.GeneratePem()
			server := httptest.NewServer(testHandlers())
			client := Client{ApiUri: server.URL, Pem: pm}
			_, err := client.PairWithCode("abcdefg")
			Expect(err).To(MatchError("407: This error is fake"))
		})
		It("Makes a POST to the tokens endpoint", func() {
			pm := ku.GeneratePem()
			server := httptest.NewServer(stubHandlers())
			client := Client{ApiUri: server.URL, Pem: pm}
			toke, _ := client.PairWithCode("abcdefg")
			Expect(toke.Facade).To(Equal("public"))
		})
	})
	//We do not have to check the validity of the currency, because the compiler will check for float values

	Describe("Creates Invoices", func() {
		It("Checks the validity of the currency", func() {
			client := new(Client)
			_, err := client.CreateInvoice(10, "USDA")
			Expect(err).To(MatchError("BitPayArgumentError: invalid currency code"))
		})
		It("makes a POST to the invoices endpoint", func() {
			server := httptest.NewServer(stubHandlers())
			var toke Token
			toke.Token = "aval"
			pm := ku.GeneratePem()
			client := Client{ApiUri: server.URL, Token: toke, Pem: pm}
			invoice, _ := client.CreateInvoice(10, "USD")
			Expect(invoice.Price).To(Equal(10.0))
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

func stubHandlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/tokens", returnTokenHandler).Methods("POST")
	r.HandleFunc("/invoices", returnInvoiceHandler).Methods("POST")
	return r
}

func createTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(407)
	w.Write([]byte("{\n  \"error\": \"This error is fake\"\n}"))
}

func returnInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{\n  \"facade\": \"pos\", \"data\": {\n \"price\": 10 \n} \n}"))
}

func returnTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{\n  \"data\": [{\n \"facade\": \"public\", \"dateCreated\": 123456789 \n}] \n}"))
}
