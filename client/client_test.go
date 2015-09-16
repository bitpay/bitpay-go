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
			Expect(toke.Facade).To(Equal("pos"))
		})
	})
	Describe("PairWithFacade", func() {
		It("Makes a POST to the tokens endpoint", func() {
			pm := ku.GeneratePem()
			server := httptest.NewServer(stubHandlers())
			client := Client{ApiUri: server.URL, Pem: pm}
			toke, _ := client.PairWithFacade("merchant")
			Expect(toke.Facade).To(Equal("pos"))
		})
	})
	//We do not have to check the validity of the currency, because the compiler will check for float values

	Describe("GetTokens", func() {
		It("Makes a GET to the tokens endpoint", func() {
			pm := ku.GeneratePem()
			server := httptest.NewServer(stubHandlers())
			client := Client{ApiUri: server.URL, Pem: pm}
			tokes, _ := client.GetTokens()
			expected_token := map[string]string{"pos": "5hQxmsrK4DStwVHCnv8stryPWBSdWH8SZpNHSoEZjZxw"}
			Expect(tokes[0]).To(Equal(expected_token))
		})
	})

	Describe("Post", func() {
		It("Sends a post request to the input endpoint", func() {
			pm := ku.GeneratePem()
			server := httptest.NewServer(stubHandlers())
			client := Client{ApiUri: server.URL, Pem: pm}
			paylo := make(map[string]string)
			response, _ := client.Post("random", paylo)
			Expect(response.StatusCode).To(Equal(209))
		})
	})

	Describe("GetToken", func() {
		It("Gets a single token as a string from a specified facade", func() {
			pm := ku.GeneratePem()
			server := httptest.NewServer(stubHandlers())
			client := Client{ApiUri: server.URL, Pem: pm}
			token, _ := client.GetToken("mercante")
			expected_token := "6Xt3pgLsSgVDrBdjHuqfJcjmGLVst2KxZYe7fqfmnUmB"
			Expect(token).To(Equal(expected_token))
		})
		It("Errors out if there is no token for the facade", func() {
			pm := ku.GeneratePem()
			server := httptest.NewServer(stubHandlers())
			client := Client{ApiUri: server.URL, Pem: pm}
			_, err := client.GetToken("mercmerc")
			expected_error := "facade not available in tokens"
			Expect(err).To(MatchError(expected_error))
		})
	})

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
	r.HandleFunc("/tokens", returnAllTokensHandler).Methods("GET")
	r.HandleFunc("/random", createRandom).Methods("POST")
	return r
}

func createTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(407)
	w.Write([]byte("{\n  \"error\": \"This error is fake\"\n}"))
}

func createRandom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(209)
	w.Write([]byte("{\n  \"no error\": \"This is not an error\"\n}"))
}

func returnInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{\n  \"facade\": \"pos\", \"data\": {\n \"price\": 10 \n} \n}"))
}

func returnTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{\n  \"data\": [{\n \"facade\": \"pos\", \"dateCreated\": 123456789 \n}] \n}"))
}

func returnAllTokensHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{\n  \"data\": [\n    {\n      \"pos\": \"5hQxmsrK4DStwVHCnv8stryPWBSdWH8SZpNHSoEZjZxw\"\n    },\n    {\n      \"pos\": \"72CFmTt5d5J9wH9R3hp69yTxtRTcA3nuDC9VurmvvUfw\"\n    },\n    {\n      \"merchant\": \"6Xt2pgLsSgVDrBdjHuqfJcjmGLVst2KxZYe7fqfmnUmB\"\n    } ,\n    {\n      \"mercante\": \"6Xt3pgLsSgVDrBdjHuqfJcjmGLVst2KxZYe7fqfmnUmB\"\n    } ,\n    {\n      \"merchant\": \"6Xt4pgLsSgVDrBdjHuqfJcjmGLVst2KxZYe7fqfmnUmB\"\n    },\n    {\n      \"mercante\": \"6Xt5pgLsSgVDrBdjHuqfJcjmGLVst2KxZYe7fqfmnUmB\"\n    }  \n  ]\n}"))
}
