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
			Expect(err).To(MatchError("BitPayArgumentError: invalid pairing code"))
		})
		It("Handles Errors Gracefully", func() {
			server := httptest.NewServer(testHandlers())
			client := Client{ApiUri: server.URL}
			_, err := client.PairWithCode("abcdefg")
			Expect(err).To(MatchError("407: This error is fake"))
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
			server := httptest.NewServer(testHandlers())
			client := Client{ApiUri: server.URL}
			invoice, _ := client.CreateInvoice(10, "USD")
			Expect(invoice["price"]).To(Equal("10"))
		})
	})
})

//handlers function
//each handler function

func testHandlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/tokens", createTokenHandler).Methods("POST")
	r.HandleFunc("/invoices", createInvoiceHandler).Methods("POST")
	return r
}

func createTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(407)
	w.Write([]byte("{\n  \"error\": \"This error is fake\"\n}"))
}

func createInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{\n  \"facade\": \"pos\", \"data\": {\n \"price\": \"10\" \n} \n}"))
}
