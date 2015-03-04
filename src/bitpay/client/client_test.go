package client_test

import (
	. "bitpay/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Context("I hate this", func() {
		It("Should be true", func() {
			Expect(ThisIsAFunction()).To(Equal(true))
		})
	})
})
