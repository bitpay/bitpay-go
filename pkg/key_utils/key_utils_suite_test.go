package key_utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKeyUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KeyUtils Suite")
}
